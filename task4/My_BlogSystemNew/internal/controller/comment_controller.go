package controller

import (
	"blog-system/internal/dto"
	"blog-system/internal/service"
	"blog-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService service.CommentService
}

func NewCommentController(commentService service.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
	}
}

// CreateComment 创建评论
// @Summary 创建评论
// @Description 对文章创建评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCommentRequest true "评论信息"
// @Success 201 {object} utils.Response{data=dto.CommentResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/comments [post]
func (c *CommentController) CreateComment(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	var req dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	response, err := c.commentService.CreateComment(&req, userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, http.StatusCreated, "评论创建成功", response)
}

// GetCommentsByPost 获取文章评论
// @Summary 获取文章评论
// @Description 获取指定文章的所有评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param postId path int true "文章ID"
// @Success 200 {object} utils.Response{data=[]dto.CommentResponse}
// @Router /api/posts/{postId}/comments [get]
func (c *CommentController) GetCommentsByPost(ctx *gin.Context) {
	postID, err := strconv.ParseUint(ctx.Param("postId"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	comments, err := c.commentService.GetCommentsByPostID(uint(postID))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "获取成功", comments)
}

// GetUserComments 获取用户评论
// @Summary 获取用户评论
// @Description 获取当前用户的评论列表
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]dto.CommentResponse}
// @Router /api/comments/my [get]
func (c *CommentController) GetUserComments(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	comments, err := c.commentService.GetCommentsByUserID(userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "获取成功", comments)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Description 删除评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/comments/{id} [delete]
func (c *CommentController) DeleteComment(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的评论ID")
		return
	}

	err = c.commentService.DeleteComment(uint(id), userID.(uint))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "评论不存在" {
			status = http.StatusNotFound
		} else if err.Error() == "无权删除此评论" {
			status = http.StatusForbidden
		}
		utils.Error(ctx, status, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "评论删除成功", nil)
}