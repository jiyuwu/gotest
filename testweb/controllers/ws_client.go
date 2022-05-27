package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jiyuwu/gotest/testweb/vo"
)

const (
	defaultAppId = 101 // 默认平台Id
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

var (
	clientManager = NewClientManager()                    // 管理者
	appIds        = []uint32{defaultAppId, 102, 103, 104} // 全部的平台

	serverIp   string
	serverPort string
)

func GetServer() (server *vo.Server) {
	server = vo.NewServer(serverIp, serverPort)

	return
}

func IsLocal(server *vo.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}

	return
}

// Client is a websocket client
type Client struct {
	Addr          string // 客户端地址
	Token         string // 用户token，用户登录以后才有
	UserId        string // 用户Id 登录以后才有
	LoginTime     uint64 // 登录时间 登录以后才有
	FirstTime     uint64 // 首次连接时间
	HeartbeatTime uint64 // 用户上次心跳时间
	AppId         uint32 // 登录的平台Id app/web/ios
	Socket        *websocket.Conn
	Send          chan []byte // 待发送的数据
}

// 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// 读取客户端数据
func (c *Client) SendMsg(msg []byte) {

	if c == nil {

		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}

// 用户登录
func (c *Client) Login(appId uint32, token string, userId string, loginTime uint64) {
	c.AppId = appId
	c.Token = token
	c.UserId = userId
	c.LoginTime = loginTime
	// 登录成功=心跳一次
	c.Heartbeat(loginTime)
}

// 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}

// 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := clientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.LoginTime, client.HeartbeatTime)
			client.Socket.Close()
			clientManager.Unregister <- client
		}
	}
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case conn := <-clientManager.Register:
			log.Printf("新用户加入:%v", conn.Token)
			clientManager.AddClients(conn)

			jsonMessage, _ := json.Marshal(vo.NewResponseHead("", "Register", 123, "", "新用户加入"))
			conn.Send <- jsonMessage
		case conn := <-clientManager.Unregister:
			log.Printf("用户离开:%v", conn.Token)
			clientManager.DelClients(conn)

			jsonMessage, _ := json.Marshal(vo.NewResponseHead("", "Unregister", 123, "", "用户离开"))
			conn.Send <- jsonMessage
		}
	}
}

func (c *Client) Read() {
	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			clientManager.Unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("读取到客户端的信息:%s", string(message))
		//clientManager.Broadcast <- message
		ProcessData(c, message)
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
