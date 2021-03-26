package controllers

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommunityHandler 查询社区列表
// @Summary 查询社区列表接口
// @Description 查询已有帖子
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunityList
// @Router /community/:id [get]
func CommunityHandler(c *gin.Context) {
	//查询到所有的社区（community_id community_name）以切片形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
// @Summary 社区分类详情接口
// @Description 根据url中的community_id查询对应帖子详情
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamCommunityID true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponseCommunity
// @Router /community [get]
func CommunityDetailHandler(c *gin.Context) {
	//将查询的社区数据返回
	//1.获取社区id
	IDStr := c.Param("id")
	communityID, err := strconv.ParseInt(IDStr, 10, 64) //str转int
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.根据id查询
	data, err := logic.GetCommunityDetail(communityID)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}
