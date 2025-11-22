package repository

import (
	"blog-system/internal/model"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *model.Comment) error
	FindByID(id uint) (*model.Comment, error)
	FindByPostID(postID uint) ([]model.Comment, error)
	FindByUserID(userID uint) ([]model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// 说明：CommentRepository 负责评论的持久化操作，包含按文章查询与评论管理方法。

func (r *commentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").Preload("Post").First(&comment, id).Error
	return &comment, err
}

func (r *commentRepository) FindByPostID(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").
		Where("post_id = ? AND status = ?", postID, "approved").
		Where("parent_id IS NULL").
		Order("created_at ASC").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) FindByUserID(userID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("Post").Where("user_id = ?", userID).Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}
