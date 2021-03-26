package middlewares

import (
	"bluebell/controllers"
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件。这个中间件判断token，如果能运行完就说明符合要求，且token绑定到c *gin.Context中
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头 样例:Authorization:Bearer xxxxxx.xx.xx
		// 但是具体实现方式要依据你的实际业务情况决定
		//1.判断请求头重是否带auth格式的JWT token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			c.Abort()
			return
		}
		//处理token
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort() // Abort prevents pending handlers from being called. Note that this will not stop the current handler.
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userid信息保存到请求的上下文c上
		c.Set(controllers.ContextUserIDKey, mc.UserID)
		c.Next() // 后续的处理请求函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息。继续后面的函数？
	}
}
