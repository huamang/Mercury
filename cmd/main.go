package main

import (
	"Mercury/common"
	"Mercury/engine"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	common.Parse()
	engine.Scan()
	end := time.Now().Sub(start)
	fmt.Printf("[*] Spend time: %s\n", end)
}
