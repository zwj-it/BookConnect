package models

import "time"

type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"  form:"community_id" `
	Name         string    `json:"name" db:"community_name"  form:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction" form:"introduction"` //omitempty字段为空不展示
	CreateTime   time.Time `json:"create_time" db:"create_time" form:"create_time"`
}
