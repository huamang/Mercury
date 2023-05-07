package common

import "os"

// Init 初始化
func Init() {
	// flag初始化
	Parse()
	if Host == "" && HostFile == "" {
		os.Exit(0)
		return
	}
	// logger初始化
	InitLogger()
}
