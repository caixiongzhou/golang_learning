package controller

import (
	"blog-system/internal/dto"
	"blog-system/internal/service"
	"blog-system/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 201 {object} utils.Response{data=dto.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	response, err := c.authService.Register(&req)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, http.StatusCreated, "注册成功", response)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} utils.Response{data=dto.AuthResponse}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	response, err := c.authService.Login(&req)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "登录成功", response)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=model.User}
// @Failure 401 {object} utils.Response
// @Router /api/auth/profile [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Error(ctx, http.StatusUnauthorized, "未授权访问")
		return
	}

	user, err := c.authService.GetUserProfile(userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(ctx, http.StatusOK, "获取成功", user)
}