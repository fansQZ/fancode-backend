package po

import "gorm.io/gorm"

// ProblemBank  题库
type ProblemBank struct {
	gorm.Model
	Name        string     `gorm:"column:name" json:"name"`
	Icon        string     `gorm:"column:icon" json:"icon"`
	Description string     `gorm:"column:description" json:"description"`
	CreatorID   uint       `gorm:"column:creator_id" json:"creatorID"`
	Problems    []*Problem `gorm:"foreignKey:bank_id" json:"problems"`
}
