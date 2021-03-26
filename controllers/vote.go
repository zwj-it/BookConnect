package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// PostVoteController 帖子投票
// @Summary 帖子投票接口
// @Description 为帖子投票，不允许重复投票
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamVoteData true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseUser
// @Router /login [post]
func PostVoteController(c *gin.Context) {
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //做个类型断言，什么意思。应该是判断err这个错误是不是后面括号里的错误
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans)) //翻译错误信息，并把错误提示中的结构体标识去掉，较早之前的视频有讲这个翻译器
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	userID, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("getCurrentUser failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServeBusy, "不允许重复投票")
		return
	}
	ResponseSuccess(c, nil)
}
