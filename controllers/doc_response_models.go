package controllers

import "bluebell/models"

//专门放接口文档用到的model
//这样做的原因是因为我们的接口文档返回的数据格式是一致的，但是数据类型是不同的
type _ResponsePostList struct {
	Code    ResponseCode            `json:"code"`    //业务状态响应码
	Message string                  `json:"message"` //提示信息
	Data    []*models.ApiPostDetail `json:"data"`    //返回数据
}
type _ResponsePost struct {
	Code    ResponseCode `json:"code"`    //业务状态响应码
	Message string       `json:"message"` //提示信息
	//Data    []*models.ApiPostDetail `json:"data"`    //返回数据
}
type _ResponseCommunityList struct {
	Code    ResponseCode `json:"code"`    //业务状态响应码
	Message string       `json:"message"` //提示信息
	//Data    []*models.ApiPostDetail `json:"data"`    //返回数据
}
type _ResponseCommunity struct {
	Code    ResponseCode              `json:"code"`    //业务状态响应码
	Message string                    `json:"message"` //提示信息
	Data    []*models.CommunityDetail `json:"data"`    //返回数据
}
type _ResponseUser struct {
	Code    ResponseCode `json:"code"`    //业务状态响应码
	Message string       `json:"message"` //提示信息
	//Data    []*models.ApiPostDetail `json:"data"`    //返回数据
}
type _ResponseVote struct {
	Code    ResponseCode `json:"code"`    //业务状态响应码
	Message string       `json:"message"` //提示信息
	//Data    []*models.ApiPostDetail `json:"data"`    //返回数据
}
