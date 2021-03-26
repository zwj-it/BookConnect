package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"math"
	"time"
)

// 推荐阅读
// 基于用户投票的相关算法：http://www.ruanyifeng.com/blog/algorithm/

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200  --> 200张赞成票可以给你的帖子续一天

/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录 两次投票差值的绝对值1 +432
	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录 两次投票差值的绝对值2 +432 *2
direction=0时，有两种情况：
	1. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录 两次投票差值的绝对值1 +432
	2. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录	两次投票差值的绝对值1 -432
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录	两次投票差值的绝对值1 -432
	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录	两次投票差值的绝对值2 -432 *2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票432分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	posttime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-posttime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	} //超出限制时间不能投票。

	//2和3应在同一个事物中，同时成功或同时失败
	//2.更新帖子分数
	//2.1查询当前用户给当前帖子的投票记录（1赞成，-1反对，0没投过）
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPf+postID), userID).Val()
	//如果当前投票值和之前的相同，则提示不允许重复投票
	zap.L().Debug("VoteForPost", zap.Any("ov", ov), zap.Any("values", value))
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64 //判断当前投票是给帖子加分还是减分。
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value)
	pipeline := client.TxPipeline()
	//计算两次投票的差值
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID) //redis相关，不懂

	//3.记录用户为该帖子投票的数据
	if value == 0 { //取消投票。
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPf+postID), userID)
	} else { //记录投票数据
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPf+postID), redis.Z{
			Score:  value, //当前用户对帖子的投票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
