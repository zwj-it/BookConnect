package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code": 10001,	//程序中的错误
	"msg": xx, 	  	//提示信息
	"data": {},   	//数据
}
*/
type ResponseData struct {
	Code ResponseCode `json:"code"`
	Msg  interface{}  `json:"msg"` //定义空接口是因为msg的数据类型有很多种，这样定义限制宽松
	Data interface{}  `json:"data,omitempty"`
}

//错误返回响应
func ResponseError(c *gin.Context, code ResponseCode) {
	////上面的结构体也可以定义为gin.H,gin.H本身是个map[string]interface{}类型
	//gin.H{
	//	"code":"xx",
	//	"msg":"xx",
	//	"data":"xx",
	//}
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

//返回带数据的响应
func ResponseErrorWithMsg(c *gin.Context, code ResponseCode, data interface{}) {
	////上面的结构体也可以定义为gin.H,gin.H本身是个map[string]interface{}类型
	//gin.H{
	//	"code":"xx",
	//	"msg":"xx",
	//	"data":"xx",
	//}
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: data,
	})
}

//返回成功响应,返回成功不用像返回错误那样定义一个code ResponseCode是因为返回有多种提示，成功只有一种
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
