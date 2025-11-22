package controller

import (
	"blog-system/internal/dto"
	"blog-system/internal/service"
	"blog-system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

// CreatePost 创建文章
// @Summary 创建文章
// @Description 创建新的博客文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePostRequest true "文章信息"
// @Success 201 {object} utils.Response{data=dto.PostResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/posts [post]
func (c *PostController) CreatePost(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	var req dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	response, err := c.postService.CreatePost(&req, userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, http.StatusCreated, "文章创建成功", response)
}

// GetPost 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response{data=dto.PostResponse}
// @Failure 404 {object} utils.Response
// @Router /api/posts/{id} [get]
func (c *PostController) GetPost(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	response, err := c.postService.GetPostByID(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "获取成功", response)
}

// GetPosts 获取文章列表
// @Summary 获取文章列表
// @Description 获取分页的文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]dto.PostResponse}
// @Router /api/posts [get]
func (c *PostController) GetPosts(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	posts, total, err := c.postService.GetAllPosts(page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"hasNext":  int64(page*pageSize) < total,
		},
	}

	utils.Success(ctx, http.StatusOK, "获取成功", response)
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Description 更新文章信息
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Param request body dto.UpdatePostRequest true "文章信息"
// @Success 200 {object} utils.Response{data=dto.PostResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/posts/{id} [put]
func (c *PostController) UpdatePost(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var req dto.UpdatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	response, err := c.postService.UpdatePost(uint(id), &req, userID.(uint))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "文章不存在" {
			status = http.StatusNotFound
		} else if err.Error() == "无权修改此文章" {
			status = http.StatusForbidden
		}
		utils.Error(ctx, status, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "文章更新成功", response)
}

// DeletePost 删除文章
// @Summary 删除文章
// @Description 删除文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/posts/{id} [delete]
func (c *PostController) DeletePost(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	err = c.postService.DeletePost(uint(id), userID.(uint))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "文章不存在" {
			status = http.StatusNotFound
		} else if err.Error() == "无权删除此文章" {
			status = http.StatusForbidden
		}
		utils.Error(ctx, status, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "文章删除成功", nil)
}

// GetUserPosts 获取用户文章
// @Summary 获取用户文章
// @Description 获取指定用户的文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response{data=[]dto.PostResponse}
// @Router /api/posts/my [get]
func (c *PostController) GetUserPosts(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	posts, total, err := c.postService.GetUserPosts(userID.(uint), page, pageSize)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"posts": posts,
		"pagination": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"hasNext":  int64(page*pageSize) < total,
		},
	}

	utils.Success(ctx, http.StatusOK, "获取成功", response)
}

// GetUserPost 获取用户文章详情
// @Summary 获取用户文章详情
// @Description 根据ID获取指定用户的文章详情
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "文章ID"
// @Success 200 {object} utils.Response{data=dto.PostResponse}
// @Failure 404 {object} utils.Response
// @Router /api/posts/user/{id} [get]
func (c *PostController) GetUserPost(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "无效的文章ID")
		return
	}

	// 需要先验证文章属于当前用户
	response, err := c.postService.GetPostByIDAndUser(uint(id), userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "获取成功", response)
}
