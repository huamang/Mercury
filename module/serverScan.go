package module

import (
	"Mercury/common"
	"Mercury/plugins/redis"
)

func ServerScan(taskList []ScanTask) {
	taskChan := make(chan ScanTask, len(taskList))
	for _, task := range taskList {
		wg.Add(1)
		taskChan <- task
	}
	// 优先进行特殊服务识别
	for i := 0; i < common.Thread; i++ {
		for task := range taskChan {
			if task.Server == "" {
				switch task.Port {
				case "6379":
					redis.Check(task.Addr)
				}
			}
		}
	}
}
