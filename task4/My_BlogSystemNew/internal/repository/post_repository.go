package repository

import (
	"blog-system/internal/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *model.Post) error
	FindByID(id uint) (*model.Post, error)
	FindAll(page, pageSize int) ([]model.Post, int64, error)
	FindByUserID(userID uint, page, pageSize int) ([]model.Post, int64, error)
	Update(post *model.Post) error
	Delete(id uint) error
	IncrementViewCount(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// 说明：PostRepository 管理文章的数据访问，所有与文章相关的数据库查询
// 都应该在此处封装，便于后续替换实现或添加缓存层。

func (r *postRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	return &post, err
}

func (r *postRepository) FindAll(page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	r.db.Model(&model.Post{}).Where("status = ?", "published").Count(&total)

	// 获取分页数据
	err := r.db.Preload("User").
		Where("status = ?", "published").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offset := (page - 1) * pageSize

	r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

func (r *postRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}
