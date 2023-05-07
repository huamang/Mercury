package common

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func ParseIP(host string, filename string) (hosts []string, err error) {
	// 从文件中读取
	if filename != "" {
		var fileHost []string
		fileHost = ReadIPFile(filename)
		hosts = append(hosts, fileHost...)
		return hosts, nil
	}
	// 从命令行参数中读取
	hosts = ParseIPs(host)
	// 去重返回
	return RemoveDuplicate(hosts), nil
}

func ParseIPs(ip string) (hosts []string) {
	switch {
	// 解析CIDR
	case strings.Contains(ip, "/"):
		return ParseIPCIDR(ip)

	// 192.168.1.1-192.168.1.100
	case strings.Contains(ip, "-"):
		return parseIP1(ip)

	// 192.168.1.1,192.168.1.2
	case strings.Contains(ip, ","):
		return parseIP2(ip)

	// 处理单个ip
	default:
		testIP := net.ParseIP(ip)
		if testIP == nil {
			return nil
		}
		return []string{ip}
	}
}

// ParseIPCIDR IP address with CIDR notation
func ParseIPCIDR(host string) (hosts []string) {
	ip, ipNet, err := net.ParseCIDR(host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"host": host,
			"err":  err,
		}).Fatal("Invalid CIDR address")
	}

	for ip = ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		hosts = append(hosts, ip.String())
	}
	return
}

// parseIP1 192.168.1.1-192.168.1.100 or 192.168.1.1-255
func parseIP1(host string) (hosts []string) {
	IPRange := strings.Split(host, "-")
	testIP := net.ParseIP(IPRange[0])
	if len(IPRange[1]) < 4 {
		Range, err := strconv.Atoi(IPRange[1])
		if testIP == nil || Range > 255 || err != nil {
			return nil
		}
		SplitIP := strings.Split(IPRange[0], ".")
		ip1, err1 := strconv.Atoi(SplitIP[3])
		ip2, err2 := strconv.Atoi(IPRange[1])
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return nil
		}
		for i := ip1; i <= ip2; i++ {
			hosts = append(hosts, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		SplitIP1 := strings.Split(IPRange[0], ".")
		SplitIP2 := strings.Split(IPRange[1], ".")
		if len(SplitIP1) != 4 || len(SplitIP2) != 4 {
			return nil
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(SplitIP1[i])
			ip2, err2 := strconv.Atoi(SplitIP2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return nil
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			hosts = append(hosts, ip)
		}
	}
	return hosts
}

// parseIP2 192.168.1.1,192.168.1.2
func parseIP2(host string) (hosts []string) {
	ipList := strings.Split(host, ",")
	for _, ipStr := range ipList {
		ip := net.ParseIP(strings.TrimSpace(ipStr))
		if ip == nil {
			logrus.WithFields(logrus.Fields{
				"host": host,
				"err":  fmt.Sprintf("Invalid IP address: %s", ipStr),
			}).Fatal("Invalid IP address")
		}

		hosts = append(hosts, ip.String())
	}
	return
}

// ReadIPFile 从文件中读取IP地址
func ReadIPFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var hosts []string
	for scanner.Scan() {
		ip := scanner.Text()
		hosts = append(hosts, ParseIPs(ip)...)
	}
	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return hosts
}

// Helper function to increment IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// RemoveDuplicate 去重
func RemoveDuplicate(old []string) []string {
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
