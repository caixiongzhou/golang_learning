package service

import (
	"blog-system/internal/dto"
	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/pkg/logger"
	"errors"
)

type CommentService interface {
	CreateComment(req *dto.CreateCommentRequest, userID uint) (*dto.CommentResponse, error)
	GetCommentsByPostID(postID uint) ([]dto.CommentResponse, error)
	GetCommentsByUserID(userID uint) ([]dto.CommentResponse, error)
	DeleteComment(id uint, userID uint) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	log         *logger.Logger
}

func NewCommentService(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, log *logger.Logger) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		log:         log,
	}
}

func (s *commentService) CreateComment(req *dto.CreateCommentRequest, userID uint) (*dto.CommentResponse, error) {
	// 检查文章是否存在
	_, err := s.postRepo.FindByID(req.PostID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 检查父评论是否存在（如果提供了ParentID）
	if req.ParentID != nil {
		parentComment, err := s.commentRepo.FindByID(*req.ParentID)
		if err != nil || parentComment.PostID != req.PostID {
			return nil, errors.New("父评论不存在")
		}
	}

	comment := &model.Comment{
		Content:  req.Content,
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
		Status:   "approved", // 简化处理，直接批准
	}

	if err := s.commentRepo.Create(comment); err != nil {
		s.log.Error("创建评论失败: %v", err)
		return nil, errors.New("创建评论失败")
	}

	// 重新加载以获取关联数据
	newComment, err := s.commentRepo.FindByID(comment.ID)
	if err != nil {
		return nil, errors.New("创建评论失败")
	}

	return s.convertToCommentResponse(newComment), nil
}

// 说明：CommentService 负责评论的业务处理，包括校验目标文章、保存评论并返回格式化的响应。
// - 该层可扩展审核流程、反垃圾/敏感词检查等业务逻辑。
func (s *commentService) GetCommentsByPostID(postID uint) ([]dto.CommentResponse, error) {
	comments, err := s.commentRepo.FindByPostID(postID)
	if err != nil {
		s.log.Error("获取评论失败: %v", err)
		return nil, errors.New("获取评论失败")
	}

	// 构建评论树
	commentMap := make(map[uint]*dto.CommentResponse)
	var rootComments []*dto.CommentResponse

	// 第一遍：创建所有评论的响应对象
	for i := range comments {
		response := s.convertToCommentResponse(&comments[i])
		commentMap[response.ID] = response
	}

	// 第二遍：构建树形结构
	for i := range comments {
		response := commentMap[comments[i].ID]
		if comments[i].ParentID == nil {
			rootComments = append(rootComments, response)
		} else {
			parent := commentMap[*comments[i].ParentID]
			parent.Replies = append(parent.Replies, *response)
		}
	}

	// 转换为切片
	result := make([]dto.CommentResponse, len(rootComments))
	for i, comment := range rootComments {
		result[i] = *comment
	}

	return result, nil
}

func (s *commentService) GetCommentsByUserID(userID uint) ([]dto.CommentResponse, error) {
	comments, err := s.commentRepo.FindByUserID(userID)
	if err != nil {
		s.log.Error("获取用户评论失败: %v", err)
		return nil, errors.New("获取用户评论失败")
	}

	responses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = *s.convertToCommentResponse(&comment)
	}

	return responses, nil
}

func (s *commentService) DeleteComment(id uint, userID uint) error {
	comment, err := s.commentRepo.FindByID(id)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 检查权限
	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}

	if err := s.commentRepo.Delete(id); err != nil {
		s.log.Error("删除评论失败: %v", err)
		return errors.New("删除评论失败")
	}

	return nil
}

func (s *commentService) convertToCommentResponse(comment *model.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		Status:    comment.Status,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Author: dto.AuthorInfo{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Email:    comment.User.Email,
			Avatar:   comment.User.Avatar,
		},
		Replies: []dto.CommentResponse{},
	}
}
