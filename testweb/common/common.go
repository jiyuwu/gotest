package common

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	OK                 = 200  // Success
	NotLoggedIn        = 1000 // 未登录
	ParameterIllegal   = 1001 // 参数不合法
	UnauthorizedUserId = 1002 // 非法的用户Id
	Unauthorized       = 1003 // 未授权
	ServerError        = 1004 // 系统错误
	NotData            = 1005 // 没有数据
	ModelAddError      = 1006 // 添加错误
	ModelDeleteError   = 1007 // 删除错误
	ModelStoreError    = 1008 // 存储错误
	OperationFailure   = 1009 // 操作失败
	RoutingNotExist    = 1010 // 路由不存在
)

// 根据错误码 获取错误信息
func GetErrorMessage(code uint32, message string) string {
	var codeMessage string
	codeMap := map[uint32]string{
		OK:                 "Success",
		NotLoggedIn:        "未登录",
		ParameterIllegal:   "参数不合法",
		UnauthorizedUserId: "非法的用户Id",
		Unauthorized:       "未授权",
		NotData:            "没有数据",
		ServerError:        "系统错误",
		ModelAddError:      "添加错误",
		ModelDeleteError:   "删除错误",
		ModelStoreError:    "存储错误",
		OperationFailure:   "操作失败",
		RoutingNotExist:    "路由不存在",
	}

	if message == "" {
		if value, ok := codeMap[code]; ok {
			// 存在
			codeMessage = value
		} else {
			codeMessage = "未定义错误类型!"
		}
	} else {
		codeMessage = message
	}

	return codeMessage
}
func GetServerIp() string {
	ip, err := ExternalIP()
	if err != nil {
		return ""
	}
	return ip.String()
}
func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ip := getIpFromAddr(addr); ip != nil {
				return ip, nil
			}
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	return ip.To4()
}

func GetOrderIdTime() (orderId string) {

	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)

	return
}
