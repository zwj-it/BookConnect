package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	//查数据库，查到所有community并返回
	data, err = mysql.GetCommunityList()
	return
}
func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	//查数据库，查到所有community并返回
	data, err = mysql.GetCommunityDetailByID(id)
	return
}
