package rpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/rpcinterface/rpc_proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/runtime/protoiface"
)

type server struct {
}

func setErr(rsp protoiface.MessageV1, code uint32, message string) {

	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *rpc_proto.SendMsgRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

// 给本机用户发消息
func (s *server) SendMsg(c context.Context, req *rpc_proto.SendMsgReq) (rsp *rpc_proto.SendMsgRsp, err error) {

	fmt.Println("grpc_request 给本机用户发消息", req.String())

	rsp = &rpc_proto.SendMsgRsp{}

	if req.GetIsLocal() {
		// 不支持
		setErr(rsp, common.ParameterIllegal, "")

		return
	}

	sendResults, err := controllers.SendUserMessageLocal(req.GetAppId(), req.GetUserId(), req.GetMsg())
	if err != nil {
		fmt.Println("系统错误", err)
		setErr(rsp, common.ServerError, "")

		return rsp, nil
	}

	if !sendResults {
		fmt.Println("发送失败", err)
		setErr(rsp, common.OperationFailure, "")

		return rsp, nil
	}

	setErr(rsp, common.OK, "")

	fmt.Println("grpc_response 给本机用户发消息", rsp.String())
	return
}

func Init() {
	rpcPort := viper.GetString("app.rpcPort")
	rpcIp, _ := common.ExternalIP()
	fmt.Println("rpc server 启动", rpcIp.String(), rpcPort)

	lis, err := net.Listen("tcp", rpcIp.String()+":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc_proto.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
