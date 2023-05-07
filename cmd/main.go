package main

import (
	"Mercury/common"
	"Mercury/engine"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	// 初始化日志与配置
	common.Init()
	engine.Scan()
	end := time.Now().Sub(start)
	fmt.Printf("[*] Spend time: %s\n", end)
}
