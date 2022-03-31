package vo

type LoginReq struct {
	UserName string `json:"userName" form:"userName" binding:"required,gt=0"`
	Password string `json:"password" form:"password" `
}

type LoginResponse struct {
	ModelPrefix
	Token string `json:"token,omitempty"`
}
