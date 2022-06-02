package task

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/jiyuwu/gotest/testweb/controllers"
)

func Init() {
	Timer(3*time.Second, 30*time.Second, cleanConnection, "", nil, nil)

}

// 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ClearTimeoutConnections stop", r, string(debug.Stack()))
		}
	}()

	//logs.Info(fmt.Sprintf("定时任务，清理超时连接%v", param))
	controllers.ClearTimeoutConnections()

	return
}
