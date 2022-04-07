package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/gin-gonic/gin"
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

// Client is a websocket client
type Client struct {
	Addr          string // 客户端地址
	Token         string // 用户token，用户登录以后才有
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

// 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Token       chan *string       // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Token:      make(chan *string, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
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
func (c *Client) Login(appId uint32, userId string, loginTime uint64) {
	c.AppId = appId
	c.Token = userId
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

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case conn := <-clientManager.Register:
			log.Printf("新用户加入:%v", conn.Token)
			clientManager.AddClients(conn)

			jsonMessage, _ := json.Marshal(vo.NewResponseHead("", "", 123, "", "新用户加入"))
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

func WsHandler(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	client := &Client{ //连接信息
		Socket: conn,
		Send:   make(chan []byte),
	}
	clientManager.Register <- client

	go client.Read()

	go client.Write()
}
