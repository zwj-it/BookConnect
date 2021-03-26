package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	//3.ZREVRANGE 按分数从大到小的顺序查询
	return client.ZRevRange(key, start, end).Result()
}

func CreatePost(postID, communityID int64) error {
	//用到redis事务，要不同时成功要不就失败，回顾一下之前讲的redis
	pipeline := client.TxPipeline()
	//创建帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		float64(time.Now().Unix()),
		postID,
	})
	//创建帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		float64(time.Now().Unix()),
		postID,
	})
	//补充： 把帖子id加到社区的set里
	cKey := getRedisKey(KeyCommunitySetPf + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

// GetPostIDsInOrder 按Order查询帖子id，返回id列表
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id
	//1.根据用户请求中携带的key来确定要查询的redis key（根据时间或者根据分数）
	key := getRedisKey(KeyPostTimeZSet) //默认根据时间
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//zap.L().Debug("GetPostIDsInOrder", zap.Any("key", key))
	//2.确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的分数数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//不用pipeline，对每个key都查询一次
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	//查找getRedisKey(KeyPostScoreZSet+id)分数是1的元素数量，即统计每篇帖子的赞成票数量
	//	v := client.ZCount(getRedisKey(KeyPostScoreZSet+id), "1", "1").Val()
	//	data = append(data, v)
	//}
	//用pipeline,可以把一次请求都算完一起返回，减少RTT（网络通信）
	data = make([]int64, 0, len(ids))
	pipline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostScoreZSet + id)
		pipline.ZCount(key, "1", "1")
	}
	cmders, err := pipline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询id列表,在第67课讲的，redis都不懂。写完看懂redis再回来理解
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet) //默认根据时间
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore 把分区的帖子set与帖子的分数zset取!交集!，生成一个新的zset

	//使用这个新的zset 按之前的做法取数据
	//社区的key
	commuKey := getRedisKey(KeyCommunitySetPf + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	//判断这个值存不存在 这里的意思是，第一次查询一个key的时候是不存在的，需要求两个set的集合
	//而在第二次查询同一个key（间隔不超过60s时），由于有第一次的查询，redis缓存了这个key，所以可以直接去查询。
	//就是通过这个方式来减少zinterstore执行次数
	if client.Exists(key).Val() < 1 {
		//不存在，需要计算 ???
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, commuKey, orderKey) //? zinterstore计算  交集
		pipeline.Expire(key, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在则直接根据key去查询ids
	return getIDsFromKey(key, p.Page, p.Size)
}
