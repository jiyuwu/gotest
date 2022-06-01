package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/jiyuwu/gotest/testweb/cache"
	"github.com/jiyuwu/gotest/testweb/rpcinterface/rpc_client"
)

// 给本机用户发送消息
func SendUserMessageLocal(appId uint32, userId string, data string) (sendResults bool, err error) {

	client := GetUserClient(appId, userId)
	if client == nil {
		err = errors.New("用户不在线")

		return
	}

	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true

	return
}

// 给其它服务器用户发消息
func SendOtherUserMessage(appId uint32, userId string, msgId, cmd, message string) (sendResults bool, err error) {
	sendResults = true

	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给其它服务器用户发消息", err)
		return
	}

	for _, server := range servers {
		if !IsLocal(server) && server.Ip != "" {
			fmt.Println("给其它服务器用户发消息", server, msgId, appId, userId, cmd, "json", message)
			rpc_client.SendMsg(server, msgId, appId, userId, cmd, "json", message)
		}
	}
	return
}
