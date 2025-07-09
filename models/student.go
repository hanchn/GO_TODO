package models

import (
	"time"
	"gorm.io/gorm"
)

// Student 学生模型
type Student struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:100" binding:"required"`
	Age       int            `json:"age" gorm:"not null" binding:"required,min=1,max=150"`
	Gender    string         `json:"gender" gorm:"not null;size:10" binding:"required,oneof=男 女"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100" binding:"required,email"`
	Phone     string         `json:"phone" gorm:"size:20" binding:"required"`
	Major     string         `json:"major" gorm:"size:100" binding:"required"`
	Grade     string         `json:"grade" gorm:"size:50" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Student) TableName() string {
	return "students"
}