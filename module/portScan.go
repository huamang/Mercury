package module

import (
	"Mercury/common"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	ports        []int
	AliveAddress []string
	lock         sync.Mutex
)

type Addr struct {
	ip   string
	port int
}

func PortScan(hosts []string) []string {
	// 若指定port，则直接使用
	if common.Port != "" {
		ports = common.ParsePort(common.Port)
	} else {
		// 若未指定port，则使用默认端口
		ports = common.DefaultPorts
	}

	// 创建channel
	Addrs := make(chan Addr, len(hosts)*len(ports))
	results := make(chan string, len(hosts)*len(ports))
	defer func() {
		close(Addrs)
		close(results)
	}()

	// 添加扫描目标进入channel
	for _, port := range ports {
		for _, host := range hosts {
			wg.Add(1)
			Addrs <- Addr{host, port}
		}
	}

	// 指定线程数
	for i := 0; i < common.Thread; i++ {
		go func() {
			for host := range Addrs {
				ScanPort(host)
				wg.Done()
			}
		}()
	}
	wg.Wait()

	return AliveAddress
}

// ScanPort 执行tcp端口扫描
func ScanPort(addr Addr) {
	target := fmt.Sprintf("%s:%d", addr.ip, addr.port)
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err == nil || conn != nil {
		// 互斥锁
		lock.Lock()
		AliveAddress = append(AliveAddress, target)
		lock.Unlock()
	}
}
