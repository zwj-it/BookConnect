package controllers

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/settings"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建帖子
// @Summary 创建帖子接口
// @Description 根据前端传来的数据创建帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.Post true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePost
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	//1.获取参数/校验。记住，效验参数是用到gin框架里内置的validator。在定义结构体的时候用binding定义好数据的要求
	post := new(models.Post)
	if err := c.ShouldBindJSON(post); err != nil {
		zap.L().Debug("c.ShouldBindJSON(post) err", zap.Any("err", err))
		zap.L().Error("CreatePost with invalid param failed")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//1.5获得发出当前请求的用户id
	userID, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorID = userID
	//2.创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情接口
// @Description 验证用户id是否正确并
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostID true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /post/:id [get]
func GetPostDetailHandler(c *gin.Context) {
	//1.获取用户id，效验id是否正确
	pStr := c.Param("id") //路由里是id
	//pid, err := strconv.Atoi(pStr) //应该也行
	pid, err := strconv.ParseInt(pStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//3.从数据库中获取帖子详情
	data, err := logic.GetPostDetailbyID(pid)
	//4.返回响应
	if err != nil {
		zap.L().Error("logic.GetPostDetail err", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表
// @Summary 原版帖子列表接口
// @Description 根据url绑定的Page和Size查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	//获取分页参数
	page, size := getPageInfo(c)
	//获取list
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("GetPostListHandler failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
	return
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区/按时间/分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList true "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// 根据前端传来的参数动态的获取帖子列表
	//按创建时间或按照f分数排序
	//1.获取参数
	//2.取redis查询id列表
	//3.拿到id列表，根据id去数据库查询帖子详情信息
	// GET请求参数(query string)：/api/v1/posts2?page=1&size=10&order=time
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:        settings.Conf.Page,
		Size:        settings.Conf.Size,
		Order:       models.OrderTime,
		CommunityID: -1,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	//获取list
	data, err := logic.GetPostListNewest(p) //将查询所有帖子和按社区分类查询帖子合二为一
	if err != nil {
		zap.L().Error("GetPostListHandler failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
	return
}

//合并到上面了
//// GetPostListByCommunityHandler 根据社区查询帖子列表
//func GetPostListByCommunityHandler(c *gin.Context) {
//	//初始化结构体时指定初始参数
//	p := &models.ParamPostCommunityList{
//		ParamPostList: &models.ParamPostList{
//			Page:  settings.Conf.Page,
//			Size:  settings.Conf.Size,
//			Order: models.OrderTime,
//		},
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetPostListByCommunityHandler failed", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//	//c.ShouldBind()  根据请求的数据类型选择相应的方法去获取数据
//	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
//	//获取list
//	data, err := logic.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("GetCommunityPostList failed", zap.Error(err))
//		ResponseError(c, CodeServeBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//	return
//}
