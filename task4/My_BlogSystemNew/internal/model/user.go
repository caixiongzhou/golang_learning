package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email     string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	FirstName string    `gorm:"size:50" json:"first_name"`
	LastName  string    `gorm:"size:50" json:"last_name"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Bio       string    `gorm:"type:text" json:"bio"`
	Status    string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// 说明：User 实体代表系统中的用户数据结构。
// - Password 字段应存储经过 bcrypt 等算法哈希后的密文，返回给客户端时应过滤掉。
// - Posts 与 Comments 是 GORM 的关联字段，便于预加载作者相关的文章和评论信息。