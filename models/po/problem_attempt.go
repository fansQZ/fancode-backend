package po

import "gorm.io/gorm"

// ProblemAttempt  用户在一道题目中的做题情况
type ProblemAttempt struct {
	gorm.Model
	ProblemID       uint `gorm:"column:problem_id"`
	UserID          uint `gorm:"column:user_id"`
	SubmissionCount int  `gorm:"column:submission_count"`
	SuccessCount    int  `gorm:"column:success_count"`
	ErrCount        int  `gorm:"column:err_count"`
	// 最近一次的代码
	Code     string `gorm:"column:code"`
	Language string `gorm:"column:language"`
	// 0 未开始，1进行中 2 提交成功
	Status int `gorm:"column:status"`
}
