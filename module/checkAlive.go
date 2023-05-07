package module

import (
	"Mercury/common"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	resHosts []string
	wg       sync.WaitGroup
	OS       = runtime.GOOS
)

func CheckLive(hosts []string) []string {
	if common.Icmp {
		icmpCheck(hosts)
	} else {
		pingCheck(hosts)
	}
	// 去重后返回
	return RemoveDuplicate(resHosts)
}

// ping模式探活
func pingCheck(hosts []string) {
	// 限制并发数50
	limiter := make(chan struct{}, 50)
	for _, host := range hosts {
		wg.Add(1)
		limiter <- struct{}{}
		go func(host string) {
			if ExecCommandPing(host) {
				resHosts = append(resHosts, host)
				logrus.WithFields(logrus.Fields{
					"target": host,
					"status": "Alive",
				}).Info("Ping mode check alive")
			}
			<-limiter
			wg.Done()
		}(host)
	}
	wg.Wait()
}

// icmp模式探活
func icmpCheck(hosts []string) {
	//	优先尝试监听式
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	defer conn.Close()
	if err == nil {
		RunIcmp1(hosts, conn)
	} else {
		// 监听失败则ping探测
		fmt.Println("ICMP监听失败，可能需要更高权限，尝试ping探测")
		pingCheck(hosts)
	}
}

// RunIcmp1 监听模式
func RunIcmp1(hosts []string, conn *icmp.PacketConn) {
	// 构造 ICMP Echo 请求
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, // ICMP Echo 请求类型为 8
		Code: 0,
		Body: &icmp.Echo{
			ID:  os.Getpid() & 0xffff,
			Seq: 1,
		},
	}
	// 序列化 ICMP Echo 请求消息
	msgByte, err := msg.Marshal(nil)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	// 协程接收 ICMP Echo 响应
	fin := false
	go func() {
		for {
			if fin == true {
				return
			}
			resMsg := make([]byte, 100)
			// 接收 ICMP Echo 响应，如果sourceIP不为空则表示主机存活
			_, sourceIP, _ := conn.ReadFrom(resMsg)
			if sourceIP != nil {
				resHosts = append(resHosts, sourceIP.String())
				logrus.WithFields(logrus.Fields{
					"target": sourceIP.String(),
					"status": "Alive",
				}).Info("ICMP mode check alive")
			}
		}
	}()

	// 发送 ICMP Echo 请求到指定 IP
	for _, host := range hosts {
		dstAddr, err := net.ResolveIPAddr("ip4", host)
		if err != nil {
			fmt.Println("ResolveIPAddr error:", err)
			return
		}
		_, err = conn.WriteTo(msgByte, dstAddr)
		if err != nil {
			fmt.Println("WriteTo error:", err)
			return
		}
	}
	// 等待接收协程结束
	start := time.Now()
	for {
		if len(resHosts) == len(hosts) {
			break
		}
		since := time.Now().Sub(start)
		var wait time.Duration
		switch {
		case len(hosts) <= 256:
			wait = time.Second * 3
		default:
			wait = time.Second * 6
		}
		if since > wait {
			break
		}
	}

	fin = true
}

// ExecCommandPing 根据不同操作系统执行ping命令
func ExecCommandPing(host string) bool {
	if OS == "linux" {
		cmd := exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+host+">/dev/null && echo true || echo false")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "true") {
			return true
		}
	} else if OS == "windows" {
		cmd := exec.Command("cmd", "/c", "ping -n 1 -w 1 "+host+" && echo true || echo false")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "true") {
			return true
		}
	} else if OS == "darwin" {
		cmd := exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+host+" >/dev/null && echo true || echo false")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Run()
		if strings.Contains(out.String(), "true") {
			return true
		}
	} else {
		logrus.Fatal("不支持的操作系统")
	}
	return false
}

// RemoveDuplicate 去重
func RemoveDuplicate(old []string) []string {
	var result []string
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
