package dto

import (
	"FanCode/models/po"
	"time"
)

type ActivityItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// AccountInfo
// 和userInfo类似，但是比userInfo的数据多一些
type AccountInfo struct {
	Avatar       string `json:"avatar"`
	LoginName    string `json:"loginName"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Introduction string `json:"introduction"`
	Sex          int    `json:"sex"`
	BirthDay     string `json:"birthDay"`
	CodingAge    int    `json:"codingAge"`
}

func NewAccountInfo(user *po.SysUser) *AccountInfo {
	return &AccountInfo{
		Avatar:       user.Avatar,
		LoginName:    user.LoginName,
		UserName:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		Introduction: user.Introduction,
		BirthDay:     user.BirthDay.Format("2006-01-02"),
		Sex:          user.Sex,
		CodingAge:    time.Now().Year() - user.CreatedAt.Year(),
	}
}
