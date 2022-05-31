package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jiyuwu/gotest/testweb/cache"
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

	// 验证参数
	request := &vo.Login{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)

		return
	}

	// 验证token正确性
	if request.Token == "" || len(request.Token) < 10 {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)
		return
	}
	if token, err := cache.GetUserTokenInfo(vo.GetUserKey(request.AppId, request.UserId)); err != nil || token != request.Token {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)
		return
	}

	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.AppId, client.UserId, seq)
		code = common.OperationFailure

		return
	}

	// 设置在线缓存
	userOnline := vo.UserLogin(serverIp, serverPort, request.AppId, request.UserId, client.Addr, currentTime)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)

		return
	}

	// 执行登录channel
	login := &login{
		AppId:  request.AppId,
		UserId: request.UserId,
		Token:  request.Token,
		Client: client,
	}
	clientManager.Login <- login
	return
}

// 向所有人发消息
func SendAllMsgController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	//currentTime := uint64(time.Now().Unix())

	// 验证参数
	request := &vo.Msg{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("群聊 解析数据失败", seq, err)
		return
	}

	if !client.IsLogin() {
		fmt.Println("用户登录 用户未登录", client.AppId, client.UserId, seq)
		code = common.OperationFailure
		return
	}

	// 执行登录channel
	broadcast := &broadcast{
		AppId:   request.AppId,
		UserId:  request.UserId,
		GroupId: request.GroupId,
		Msg:     message,
		Client:  client,
	}
	clientManager.Broadcast <- broadcast

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
	client.Heartbeat(currentTime)

	return
}

// 请求升级长连接
func WsHandler(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), conn, currentTime)
	log.Printf("客户端RemoteAddr信息:%s", conn.RemoteAddr().String())
	go client.Read()

	go client.Write()
	clientManager.Register <- client
}
