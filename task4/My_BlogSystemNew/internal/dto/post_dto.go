package dto

import "time"

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
	Excerpt string `json:"excerpt"`
	Status  string `json:"status" binding:"oneof=draft published archived"`
	Tags    string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"min=1,max=200"`
	Content string `json:"content" binding:"min=1"`
	Excerpt string `json:"excerpt"`
	Status  string `json:"status" binding:"oneof=draft published archived"`
	Tags    string `json:"tags"`
}

type PostResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	ViewCount   int        `json:"view_count"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	UserID      uint       `json:"user_id"`
	Tags        string     `json:"tags,omitempty"`
	Author      AuthorInfo `json:"author"`
}

type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}