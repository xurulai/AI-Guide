package model

import "time"

// User 用户模型
type User struct {
	UserID    string      `json:"user_id" gorm:"primaryKey;column:user_id"`
	Username  string      `json:"username" gorm:"column:username"`
	Email     string      `json:"email" gorm:"column:email"`
	Phone     string      `json:"phone" gorm:"column:phone"`
	Profile   UserProfile `json:"profile" gorm:"serializer:json;column:profile"`
	CreatedAt time.Time   `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time   `json:"updated_at" gorm:"column:updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
