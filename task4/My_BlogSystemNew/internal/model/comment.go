package model

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    string    `gorm:"size:20;default:pending" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PostID uint `gorm:"not null;index" json:"post_id"`
	Post   Post `gorm:"foreignKey:PostID" json:"post,omitempty"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	ParentID *uint    `gorm:"index" json:"parent_id,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

// 说明：Comment 表示文章的评论，支持嵌套回复（ParentID 与 Replies）。
// - Status 可用于人工审核流程（pending/approved/rejected）。