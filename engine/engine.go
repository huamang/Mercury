package engine

import (
	"Mercury/common"
	"Mercury/module"
	"github.com/sirupsen/logrus"
	"os"
)

func Scan() {
	// IP扫描，需要探活以及端口爆破
	hosts, err := common.ParseIP(common.Host, common.HostFile)
	if err != nil {
		return
	}

	//	存活探测
	hosts = module.CheckLive(hosts)
	// 探活模式，输出后直接结束
	if common.CheckAlive {
		logrus.Info("探活模式结束")
		os.Exit(0)
	}

	// 端口扫描
	TaskList := module.PortScan(hosts)

	// 地址服务扫描
	module.ServerScan(TaskList)

}
