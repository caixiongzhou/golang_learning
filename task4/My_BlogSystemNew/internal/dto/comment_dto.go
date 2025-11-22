package dto

import "time"

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

type CommentResponse struct {
	ID        uint       `json:"id"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	PostID    uint       `json:"post_id"`
	UserID    uint       `json:"user_id"`
	ParentID  *uint      `json:"parent_id,omitempty"`
	Author    AuthorInfo `json:"author"`
	Replies   []CommentResponse `json:"replies,omitempty"`
}