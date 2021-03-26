package mysql

import (
	"bluebell/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostbyID 根据postid查询帖子详情
func GetPostbyID(pid int64) (p *models.Post, err error) {
	p = new(models.Post) //要加这个new语句，原因是？
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
from post 
where post_id = ?`
	if err := db.Get(p, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no post in db")
			err = nil
		}
	}
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
	from post
	ORDER BY create_time
	DESC 
	limit ?,?` //查询从page开始，往后size条
	postList = make([]*models.Post, 0, 2) //不要写成make([]*models.Post, 2)，要写上容量
	if err := db.Select(&postList, sqlStr, (page-1)*size, size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no post in db")
			err = nil
		}
	}
	return
}

// GetPostListByIDs 根据id列表查询帖子
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
from post
where post_id in (?)
ORDER BY FIND_IN_SET(post_id, ?)` //FIND_IN_SET根据指定顺序对查询结果排序
	qurey, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ",")) //strings.Join把字符串用逗号拼接
	if err != nil {
		return nil, err
	}
	// https: //www.liwenzhou.com/posts/Go/sqlx/
	qurey = db.Rebind(qurey)
	err = db.Select(&postList, qurey, args...) //要加... 不知道具体的格式的怎么样！！
	return
}
