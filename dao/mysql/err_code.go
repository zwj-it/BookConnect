package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在") //这种写法不要出现在代码中,拿出去用常量定义好!
	ErrorUserNoExist     = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的用户ID")
)
