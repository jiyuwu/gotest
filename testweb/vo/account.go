package vo

import (
	"fmt"
	"time"
)

//----------------------接口数据--------------------------//
type LoginReq struct {
	AppId    uint32 `json:"appId" form:"appId"` // appId
	UserName string `json:"userName" form:"userName" binding:"required,gt=0"`
	Password string `json:"password" form:"password" `
}

type LoginResponse struct {
	ModelPrefix
	Token string `json:"token,omitempty"`
	Id    int64  `json:"userId,omitempty"`
}

// 获取用户key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)

	return
}

//----------------------缓存数据--------------------------//
// 用户在线状态
type UserOnline struct {
	AccIp         string `json:"accIp"`         // acc Ip
	AccPort       string `json:"accPort"`       // acc 端口
	AppId         uint32 `json:"appId"`         // appId
	UserId        string `json:"userId"`        // 用户Id
	ClientIp      string `json:"clientIp"`      // 客户端Ip
	ClientPort    string `json:"clientPort"`    // 客户端端口
	LoginTime     uint64 `json:"loginTime"`     // 用户上次登录时间
	HeartbeatTime uint64 `json:"heartbeatTime"` // 用户上次心跳时间
	LogOutTime    uint64 `json:"logOutTime"`    // 用户退出登录的时间
	Qua           string `json:"qua"`           // qua
	DeviceInfo    string `json:"deviceInfo"`    // 设备信息
	IsLogoff      bool   `json:"isLogoff"`      // 是否下线
}

// 用户登录
func UserLogin(accIp, accPort string, appId uint32, userId string, addr string, loginTime uint64) (userOnline *UserOnline) {

	userOnline = &UserOnline{
		AccIp:         accIp,
		AccPort:       accPort,
		AppId:         appId,
		UserId:        userId,
		ClientIp:      addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}

	return
}

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

// 用户是否在线
func (u *UserOnline) IsOnline() (online bool) {
	if u.IsLogoff {

		return
	}

	currentTime := uint64(time.Now().Unix())

	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		fmt.Println("用户是否在线 心跳超时", u.AppId, u.UserId, u.HeartbeatTime)

		return
	}

	if u.IsLogoff {
		fmt.Println("用户是否在线 用户已经下线", u.AppId, u.UserId)

		return
	}

	return true
}

// 用户是否在本台机器上
func (u *UserOnline) UserIsLocal(localIp, localPort string) (result bool) {

	if u.AccIp == localIp && u.AccPort == localPort {
		result = true

		return
	}

	return
}
