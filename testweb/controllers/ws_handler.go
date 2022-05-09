package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/vo"
)

// ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)

	data = "pong"

	return
}

// 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &vo.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 用户登录", seq, "Token", request.Token)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	if request.UserId == "" || len(request.UserId) >= 20 {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)

		return
	}
	fmt.Println("用户登录 成功", seq, client.Addr, request.UserId, currentTime)
	client.Token = request.UserId
	return
}

// 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	currentTime := uint64(time.Now().Unix())

	request := &vo.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)

		return
	}

	fmt.Println("webSocket_request 心跳接口", client.AppId, client.Token)

	client.Heartbeat(currentTime)

	return
}
