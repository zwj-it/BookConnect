package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.根据数据生成postid。因为用到雪花算法
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return err
}
func GetPostDetailbyID(id int64) (data *models.ApiPostDetail, err error) {
	//查询并把我们想返回的数据拼接
	//根据id查询数据
	post, err := mysql.GetPostbyID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostbyID failed", zap.Int64("pid", id), zap.Error(err))
		return
	}
	//根据作者id，查询用户名
	user, err := mysql.GetUserbyID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserbyID failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}
	//根据社区id，查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	//把查询到的几种数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		CommunityDetail: community,
		Post:            post,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	//去数据库中查询
	posts, err := mysql.GetPostList(page, size) //查询从page开始，往后size条
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id，查询用户名
		user, err := mysql.GetUserbyID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserbyID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id，查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//把查询到的几种数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 新版获取帖子列表
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//zap.L().Debug("GetPostList2", zap.Any("p.CommunityID", p.CommunityID))
	//2.去redis查询id列表， 按分数/时间排序查询到post_id集合
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) success but len(ids)==0")
		return
	}
	//写项目的时候可以用下面这条代替输出。
	//zap.L().Debug("redis.GetPostIDsInOrder(p)", zap.Any("ids", ids))

	//2.5根据id列表，查询各个帖子的赞成票数。
	voteData, err := redis.GetPostVoteData(ids) //这里可以改，用一些复杂的算法来算这个分数，比如QBC
	if err != nil {
		return
	}
	//3.拿到id列表，根据id去数据库查询帖子详情信息
	//返回数据还要按照我给定的id顺序返回（在查询时有通过sql内置的语句限定好顺序）/也可以手动排序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GeetPostListByIDS failed", zap.Error(err))
		return
	}
	//根据posts查询帖子详细信息及分区及作者
	for idx, post := range posts {
		//根据作者id，查询用户名
		user, err := mysql.GetUserbyID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserbyID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id，查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//把查询到的几种数据拼接
		postDetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//2.去redis查询id列表， 按分数/时间排序查询到post_id集合
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) success but len(ids)==0")
		return
	}
	//写项目的时候可以用下面这条代替输出。
	//zap.L().Debug("redis.GetPostIDsInOrder(p)", zap.Any("ids", ids))

	//2.5根据id列表，查询各个帖子的赞成票数。
	voteData, err := redis.GetPostVoteData(ids) //这里可以改，用一些复杂的算法来算这个分数，比如QBC
	if err != nil {
		return
	}
	//3.拿到id列表，根据id去数据库查询帖子详情信息
	//返回数据还要按照我给定的id顺序返回（在查询时有通过sql内置的语句限定好顺序）/也可以手动排序
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GeetPostListByIDS failed", zap.Error(err))
		return
	}
	//根据posts查询帖子详细信息及分区及作者
	for idx, post := range posts {
		//根据作者id，查询用户名
		user, err := mysql.GetUserbyID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserbyID failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		//根据社区id，查询社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		//把查询到的几种数据拼接
		postDetail := &models.ApiPostDetail{
			VoteNum:         voteData[idx],
			AuthorName:      user.Username,
			CommunityDetail: community,
			Post:            post,
		}
		data = append(data, postDetail)
	}
	return
}
func GetPostListNewest(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == -1 {
		data, err = GetPostList2(p) //查询所有帖子列表
	} else {
		data, err = GetCommunityPostList(p) //按社区类别查询帖子列表
	}
	if err != nil {
		zap.L().Error("GetPostListNewest failed", zap.Error(err))
		return nil, err
	}
	return
}
