package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//定义请求的参数结构体
//ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required" form:"username" `
	Password   string `json:"password" binding:"required" form:"password" `
	RePassword string `json:"re_password" binding:"required,eqfield=Password" form:"re_password" `
}

//ParamLogin 登陆请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required" form:"username" `
	Password string `json:"password" binding:"required" form:"password" `
}

// ParamVoteData 投票
type ParamVoteData struct {
	//UserID  可以直接从当前请求中获取
	PostID    string `json:"post_id" binding:"required" form:"post_id"`                            //帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1" form:"direction" example:"1"` //投票方向，赞成1，反对-1,取消投票0

}

// ParamPostList 获取帖子列表query string参数
type ParamPostList struct {
	Page        int64  `json:"page" form:"page" example:"1"`      // 页码
	Size        int64  `json:"size" form:"size" example:"10"`     // 每页数据量
	Order       string `json:"order" form:"order" example:"time"` // 排序依据
	CommunityID int64  `json:"community_id" form:"community_id"`  //可以为空，为空则查询所有帖子
}

////ParamPostCommunityList 按社区获取帖子列表query string参数 合并到上面了
//type ParamPostCommunityList struct {
//	*ParamPostList
//	CommunityID int64 `json:"community_id" form:"community_id"`
//}

// ParamPostID 获取帖子详情需要的pid，仅在生成文档用到
type ParamPostID struct {
	//UserID  可以直接从当前请求中获取
	PostID string `json:"post_id" binding:"required" form:"post_id"` //帖子id
}

// ParamPostID 获取帖子详情需要的pid，仅在生成文档用到
type ParamCommunityID struct {
	//UserID  可以直接从当前请求中获取
	CommunityID string `json:"community_id" binding:"required" form:"community_id"` //帖子id
}
