package model

// UserProfile 用户画像
type UserProfile struct {
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID         string `gorm:"column:user_id;uniqueIndex" json:"user_id"`
	PreferredStyle string `gorm:"column:preferred_style" json:"preferred_style"` // xiaohongshu/dongyuhui
	Age            int    `gorm:"column:age" json:"age"`
	Gender         string `gorm:"column:gender" json:"gender"`
	Interests      string `gorm:"column:interests;type:json" json:"interests"` // JSON 字符串，由应用层处理
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}
