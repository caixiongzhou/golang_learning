package service

import (
	"blog-system/internal/dto"
	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/pkg/logger"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostService interface {
	CreatePost(req *dto.CreatePostRequest, userID uint) (*dto.PostResponse, error)
	GetPostByID(id uint) (*dto.PostResponse, error)
	GetAllPosts(page, pageSize int) ([]dto.PostResponse, int64, error)
	GetUserPosts(userID uint, page, pageSize int) ([]dto.PostResponse, int64, error)
	UpdatePost(id uint, req *dto.UpdatePostRequest, userID uint) (*dto.PostResponse, error)
	DeletePost(id uint, userID uint) error
	IncrementViewCount(id uint) error
	GetPostByIDAndUser(postID uint, userID uint) (*dto.PostResponse, error)
}

type postService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	log      *logger.Logger
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository, log *logger.Logger) PostService {
	return &postService{
		postRepo: postRepo,
		userRepo: userRepo,
		log:      log,
	}
}

func (s *postService) CreatePost(req *dto.CreatePostRequest, userID uint) (*dto.PostResponse, error) {
	// 生成slug
	slug := generateSlug(req.Title)

	// 设置发布时间
	var publishedAt *time.Time
	if req.Status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	post := &model.Post{
		Title:       req.Title,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Slug:        slug,
		Status:      req.Status,
		UserID:      userID,
		PublishedAt: publishedAt,
		Tags:        req.Tags,
	}

	if err := s.postRepo.Create(post); err != nil {
		s.log.Error("创建文章失败: %v", err)
		return nil, errors.New("创建文章失败")
	}

	return s.convertToPostResponse(post), nil
}

// 说明：PostService 负责文章层面的业务逻辑，包括分页查询、作者权限校验、草稿与发布流程等。
// - Service 层应尽量保持幂等与明确的错误语义，日志记录用于排查而非控制业务流程。

func (s *postService) GetPostByID(id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 增加浏览次数
	go func() {
		if err := s.postRepo.IncrementViewCount(id); err != nil {
			s.log.Error("增加文章浏览次数失败: %v", err)
		}
	}()

	return s.convertToPostResponse(post), nil
}

func (s *postService) GetAllPosts(page, pageSize int) ([]dto.PostResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := s.postRepo.FindAll(page, pageSize)
	if err != nil {
		s.log.Error("获取文章列表失败: %v", err)
		return nil, 0, errors.New("获取文章列表失败")
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = *s.convertToPostResponse(&post)
	}

	return responses, total, nil
}

func (s *postService) GetUserPosts(userID uint, page, pageSize int) ([]dto.PostResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := s.postRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		s.log.Error("获取用户文章失败: %v", err)
		return nil, 0, errors.New("获取用户文章失败")
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = *s.convertToPostResponse(&post)
	}

	return responses, total, nil
}

func (s *postService) GetPostByIDAndUser(postID uint, userID uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 验证文章属于当前用户
	if post.UserID != userID {
		return nil, errors.New("文章不存在")
	}

	return s.convertToPostResponse(post), nil
}

func (s *postService) UpdatePost(id uint, req *dto.UpdatePostRequest, userID uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 检查权限
	if post.UserID != userID {
		return nil, errors.New("无权修改此文章")
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.Status != "" {
		post.Status = req.Status
		// 如果状态变为published且尚未发布，设置发布时间
		if req.Status == "published" && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}
	if req.Tags != "" {
		post.Tags = req.Tags
	}

	if err := s.postRepo.Update(post); err != nil {
		s.log.Error("更新文章失败: %v", err)
		return nil, errors.New("更新文章失败")
	}

	return s.convertToPostResponse(post), nil
}

func (s *postService) DeletePost(id uint, userID uint) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}

	// 检查权限
	if post.UserID != userID {
		return errors.New("无权删除此文章")
	}

	if err := s.postRepo.Delete(id); err != nil {
		s.log.Error("删除文章失败: %v", err)
		return errors.New("删除文章失败")
	}

	return nil
}

func (s *postService) IncrementViewCount(id uint) error {
	return s.postRepo.IncrementViewCount(id)
}

func (s *postService) convertToPostResponse(post *model.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Excerpt:     post.Excerpt,
		Slug:        post.Slug,
		Status:      post.Status,
		ViewCount:   post.ViewCount,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
		UserID:      post.UserID,
		Tags:        post.Tags,
		Author: dto.AuthorInfo{
			ID:       post.User.ID,
			Username: post.User.Username,
			Email:    post.User.Email,
			Avatar:   post.User.Avatar,
		},
	}
}

func generateSlug(title string) string {
	// 简单的slug生成逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")

	// 添加随机后缀避免重复
	uuid := uuid.New().String()[:8]
	return slug + "-" + uuid
}
