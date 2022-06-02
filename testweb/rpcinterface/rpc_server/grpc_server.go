package rpc_server

import (
	"context"
	"fmt"
	"net"

	"github.com/jiyuwu/gotest/testweb/common"
	"github.com/jiyuwu/gotest/testweb/controllers"
	"github.com/jiyuwu/gotest/testweb/logs"
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
	logs.Info(fmt.Sprintf("grpc_request 给本机用户发消息%v", req))

	rsp = &rpc_proto.SendMsgRsp{}

	if req.GetIsLocal() {
		// 不支持
		setErr(rsp, common.ParameterIllegal, "")

		return
	}

	sendResults, err := controllers.SendUserMessageLocal(req.GetAppId(), req.GetUserId(), req.GetMsg())
	if err != nil {
		logs.Error(fmt.Sprintf("系统错误%v", err))
		setErr(rsp, common.ServerError, "")

		return rsp, nil
	}

	if !sendResults {
		logs.Error(fmt.Sprintf("发送失败%v", err))
		setErr(rsp, common.OperationFailure, "")

		return rsp, nil
	}

	setErr(rsp, common.OK, "")

	logs.Info(fmt.Sprintf("grpc_response 给本机用户发消息%s", rsp.String()))
	return
}

func Init() {
	rpcPort := viper.GetString("app.rpcPort")
	rpcIp, _ := common.ExternalIP()
	logs.Info(fmt.Sprintf("rpc server 启动%s:%s", rpcIp.String(), rpcPort))

	lis, err := net.Listen("tcp", rpcIp.String()+":"+rpcPort)
	if err != nil {
		logs.Error(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	rpc_proto.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		logs.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}
