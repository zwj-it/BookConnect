package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册
// @Summary 用户注册接口
// @Description 注册新用户
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamSignUp true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseUser
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	//1.获取参数参数效验	//通过validator库进行参数效验
	p := new(models.ParamSignUp) //这样定义是为了使用指针传递，防止结构体过大的时候影响性能
	if err := c.ShouldBindJSON(p); err != nil {
		//参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			//统一都改用上面的提示信息返回
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		return
	}
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "注册失败",
		//})

		//因为上面返回的err有多种类型所以要判断
		if errors.Is(err, mysql.ErrorUserExist) { //这个方法判断这俩错误是否相同
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil) //这里成功,但是没有啥数据要显示,所以可以用nil.也可以用"注册成功"
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "signup success",
	//})
}

// LoginHandler 用户登录
// @Summary 用户登录接口
// @Description 注册新用户
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamLogin true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseUser
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	//1.获取参数参数效验	//通过validator库进行参数效验
	p := new(models.ParamLogin) //这样定义是为了使用指针传递，防止结构体过大的时候影响性能
	if err := c.ShouldBindJSON(p); err != nil {
		//参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			//c.JSON(http.StatusOK, gin.H{
			//	"msg": err.Error(),
			//})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": removeTopStruct(errs.Translate(trans)),
		//})
		return
	}

	//2.业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNoExist) { //因为上面返回的err有多种类型所以要判断
			ResponseError(c, CodeUserNoExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		//c.JSON(http.StatusOK, gin.H{
		//	"msg": "用户名或密码错误",
		//})
		return
	}
	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), //id值大于1<<53-1   int64类型的最大值是1<<63-1,转换成string
		"user_name": user.Username,
		"token":     user.Token,
	})
	//c"uJSON(http.StatusOK, gin.H{
	//	"msg": "login success",
	//})
}
