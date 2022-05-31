package vo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type ModelPrefix struct {
	Msg  string `json:"msg"`
	Code int32  `json:"code"`
}

/************************  请求数据  **************************/
// 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            // 消息的唯一Id
	Cmd  string      `json:"cmd"`            // 请求命令字
	Data interface{} `json:"data,omitempty"` // 数据 json
}

// 登录请求数据
type Login struct {
	Token  string `json:"token"` // 验证用户是否登录
	AppId  uint32 `json:"appId,omitempty"`
	UserId string `json:"userId,omitempty"`
}

// 群聊消息
type Msg struct {
	GroupId int64  `json:"groupId"` // 组Id
	AppId   uint32 `json:"appId,omitempty"`
	UserId  string `json:"userId,omitempty"`
	Msg     string `json:"msg,omitempty"`
}

// 心跳请求数据
type HeartBeat struct {
	UserId string `json:"userId,omitempty"`
}

/************************  响应数据  **************************/
type Head struct {
	Seq      string    `json:"seq"`      // 消息的Id
	Cmd      string    `json:"cmd"`      // 消息的cmd 动作
	Response *Response `json:"response"` // 消息体
}

type Response struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

// push 数据结构体
type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// 设置返回消息
func NewResponseHead(seq string, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Seq: seq, Cmd: cmd, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)

	return
}

func NewResponse(code uint32, codeMsg string, data interface{}) *Response {
	return &Response{Code: code, CodeMsg: codeMsg, Data: data}
}

// ip port
type Server struct {
	Ip   string `json:"ip"`   // ip
	Port string `json:"port"` // 端口
}

func NewServer(ip string, port string) *Server {

	return &Server{Ip: ip, Port: port}
}

func (s *Server) String() (str string) {
	if s == nil {
		return
	}

	str = fmt.Sprintf("%s:%s", s.Ip, s.Port)

	return
}

func StringToServer(str string) (server *Server, err error) {
	list := strings.Split(str, ":")
	if len(list) != 2 {

		return nil, errors.New("err")
	}

	server = &Server{
		Ip:   list[0],
		Port: list[1],
	}

	return
}
