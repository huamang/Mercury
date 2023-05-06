package engine

import (
	"Mercury/common"
	"Mercury/module"
	"fmt"
)

// ScanIP ip扫描，需要探活以及端口爆破
func ScanIP() {
	hosts, err := common.ParseIP(common.Host, common.HostFile)
	if err != nil {
		return
	}
	//	存活探测
	hosts = module.CheckLive(hosts)
	// 探活模式，输出后直接返回
	if common.CheckAlive {
		fmt.Println("[*] 存活主机:")
		for _, host := range hosts {
			fmt.Println(host)
		}
		return
	}
	// 端口扫描
	AliveAddr := module.PortScan(hosts)
	fmt.Println("[*] 存活地址:")
	for _, addr := range AliveAddr {
		fmt.Println(addr)
	}
}
