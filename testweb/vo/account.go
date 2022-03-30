package vo

type LoginReq struct {
	UserName string `json:"userName" form:"userName" binding:"required,gt=0"`
	Password string `json:"password" form:"password" `
}
