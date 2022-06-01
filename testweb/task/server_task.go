package task

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/jiyuwu/gotest/testweb/cache"
	"github.com/jiyuwu/gotest/testweb/controllers"
)

func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// 服务注册
func server(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务注册 stop", r, string(debug.Stack()))
		}
	}()

	server := controllers.GetServer()
	currentTime := uint64(time.Now().Unix())
	fmt.Println("定时任务，服务注册", param, server, currentTime)

	cache.SetServerInfo(server, currentTime)

	return
}

// 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("服务下线 stop", r, string(debug.Stack()))
		}
	}()

	fmt.Println("服务下线", param)

	server := controllers.GetServer()
	cache.DelServerInfo(server)

	return
}
