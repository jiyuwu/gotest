package rpc_client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/rpcinterface/rpc_proto"
	"github.com/jiyuwu/gotest/testweb/vo"
	"google.golang.org/grpc"
)

// rpc client
// 发送消息
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func SendMsg(server *vo.Server, seq string, appId uint32, userId string, cmd string, msgType string, message string) (sendMsgId string, err error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(server.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("连接失败", server.String())

		return
	}
	defer conn.Close()

	c := rpc_proto.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := rpc_proto.SendMsgReq{
		Seq:     seq,
		AppId:   appId,
		UserId:  userId,
		Cms:     cmd,
		Type:    msgType,
		Msg:     message,
		IsLocal: false,
	}
	rsp, err := c.SendMsg(ctx, &req)
	if err != nil {
		fmt.Println("发送消息", err)

		return
	}

	if rsp.GetRetCode() != common.OK {
		fmt.Println("发送消息", rsp.String())
		err = errors.New(fmt.Sprintf("发送消息失败 code:%d", rsp.GetRetCode()))

		return
	}

	sendMsgId = rsp.GetSendMsgId()
	fmt.Println("发送消息 成功:", sendMsgId)

	return
}
