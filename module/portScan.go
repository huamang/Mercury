package module

import (
	"Mercury/common"
	"Mercury/fingerprint"
	"Mercury/requests"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

var (
	ports     []string
	lock      sync.Mutex
	scanTasks []ScanTask
)

func PortScan(hosts []string) []ScanTask {
	// 若指定port，则直接使用
	if common.Port != "" {
		ports = common.ParsePort(common.Port)
	} else {
		// 若未指定port，则使用默认端口
		ports = common.DefaultPorts
	}

	// 创建channel
	Addrs := make(chan string, len(hosts)*len(ports))
	defer close(Addrs)

	// 添加扫描目标进入channel
	for _, port := range ports {
		for _, host := range hosts {
			wg.Add(1)
			Addrs <- fmt.Sprintf("%s:%s", host, port)
		}
	}

	// 指定线程数，进行端口扫描
	for i := 0; i < common.Thread; i++ {
		go func() {
			for task := range Addrs {
				scanTask := ScanPort(task)
				lock.Lock()
				scanTasks = append(scanTasks, scanTask)
				lock.Unlock()
				wg.Done()
			}
		}()
	}
	wg.Wait()

	return scanTasks
}

// ScanPort 执行tcp端口扫描
func ScanPort(addr string) ScanTask {
	conn, err := requests.DialTcp(addr)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	// 若端口连接成功
	if err == nil || conn != nil {
		conn.SetReadDeadline(time.Now().Add(time.Duration(3) * time.Second))
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		buf = buf[:n]
		banner := strings.Trim(strings.TrimSpace(string(buf)), "\n")
		// banner 不为空的情况下，进行指纹识别
		if n > 0 {
			server := fingerprint.ParseBanner(buf)
			if server != "" {
				logrus.WithFields(logrus.Fields{
					"server": server,
					"target": addr,
				}).Info("Port check success")
				return ScanTask{Addr: addr, Banner: banner, Server: server}
			} else {
				logrus.WithFields(logrus.Fields{
					"banner": banner,
					"target": addr,
				}).Info("Port check success, but can't parse the banner")
				return ScanTask{Addr: addr, Banner: banner}
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"target": addr,
			}).Debug("Port check success, but banner is empty")
			return ScanTask{Addr: addr}
		}

	}
	return ScanTask{}
}
