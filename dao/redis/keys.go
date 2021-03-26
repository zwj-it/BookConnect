package redis

//redis key

//redis key尽量使用命名空间的方式，方便查询和拆分
const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"  //ZSet（数据类型）;帖子以发帖时间作为分数
	KeyPostScoreZSet   = "post:socre" //Zset;帖子及投票的分数
	KeyPostVotedZSetPf = "post:voted" //Zset;记录用户及投票类型;参数是post_id
	KeyCommunitySetPf  = "community:" //set;保存每个社区下的帖子id
)

// getRedisKey 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
