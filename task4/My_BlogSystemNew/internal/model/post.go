package model

import (
	"time"
)

type Post struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Content     string     `gorm:"type:longtext;not null" json:"content"`
	Excerpt     string     `gorm:"type:text" json:"excerpt"`
	Slug        string     `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Status      string     `gorm:"size:20;default:draft" json:"status"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`

	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Comments []Comment `gorm:"foreignKey:PostID" json:"comments,omitempty"`
	Tags     string    `gorm:"type:text" json:"tags,omitempty"`
}

// 说明：Post 表示博客文章，包含内容、状态、作者以及关联的评论。
// - Slug 用于 SEO-friendly 的文章地址，应保证唯一。
// - Tags 存储为 JSON 字符串，便于前端解析标签数组。
