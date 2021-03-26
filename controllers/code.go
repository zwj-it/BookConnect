package controllers

type ResponseCode int64

const (
	CodeSuccess ResponseCode = 1000 + iota //iota是个0？这样有啥作用
	CodeInvalidParam
	CodeUserExist
	CodeUserNoExist
	CodeInvalidPassword
	CodeServeBusy
	CodeInvalidToken
	CodeNeedLogin
)

//这样做是为了不把真正的错误返回给用户，只返回错误提示信息，真正的错误在终端或者日志中自己看
var codeMsgMap = map[ResponseCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNoExist:     "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServeBusy:       "服务繁忙",
	CodeNeedLogin:       "需要登陆",
	CodeInvalidToken:    "无效的token",
}

func (c ResponseCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServeBusy] //根据c没查到提示信息，就返回个服务繁忙。。。
	}
	return msg
}
