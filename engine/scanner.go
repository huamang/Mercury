package engine

import (
	"Mercury/common"
	"os"
)

func Scan() {
	if common.Host == "" && common.HostFile == "" {
		os.Exit(0)
		return
	}
	ScanIP()
}
