package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/vo"
	"github.com/spf13/viper"
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

func SetServer() {
	serverIp = common.GetServerIp()
	serverPort = viper.GetString("app.rpcPort")
}

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

// 读取客户端数据
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppId, c.UserId)

	return
}

// 获取用户key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)

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
type login struct {
	AppId  uint32
	UserId string
	Token  string
	Client *Client
}

// 广播发消息
type broadcast struct {
	AppId   uint32
	UserId  string
	GroupId int64
	Msg     []byte
	Client  *Client
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

// 是否登录了
func (c *Client) IsLogin() (isLogin bool) {

	// 用户登录了
	if c.UserId != "" {
		isLogin = true

		return
	}

	return
}

// 获取用户所在的连接
func GetUserClient(appId uint32, userId string) (client *Client) {
	client = clientManager.GetUserClient(appId, userId)

	return
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
		case login := <-clientManager.Login:
			// 用户登录
			clientManager.EventLogin(login)
		case message := <-clientManager.Broadcast:
			// 其它服务器rpc消息推送
			orderId := common.GetOrderIdTime()
			_, err := SendOtherUserMessage(message.AppId, message.UserId, orderId, "sendAllMsg", string(message.Msg))
			if err != nil {
				fmt.Println("SendOtherUserMessage", err.Error())
			}
			// 本地广播事件
			clients := clientManager.GetClients()
			for conn := range clients {
				if conn != message.Client {
					conn.Send <- message.Msg
				}
			}
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
