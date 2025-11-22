package service

import (
	"blog-system/internal/dto"
	"blog-system/internal/model"
	"blog-system/internal/repository"
	"blog-system/internal/utils"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 定义认证服务接口
type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) // 用户注册
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)       // 用户登录
	GetUserProfile(userID uint) (*model.User, error)              // 获取用户资料
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository // 用户数据仓库
	jwtUtil  *utils.JWTUtil            // JWT工具
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, jwtUtil *utils.JWTUtil) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

// Register 用户注册业务逻辑
func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// 检查用户名是否已存在 - 确保唯一性
	existingUser, err := s.userRepo.FindByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在 - 确保唯一性
	existingEmail, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		return nil, errors.New("注册失败")
	}

	// 创建用户
	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		log.Printf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	// 生成JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		return nil, errors.New("注册失败")
	}

	return &dto.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// Login 用户登录业务逻辑
func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("账号已被禁用")
	}

	// 生成JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		return nil, errors.New("登录失败")
	}

	return &dto.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

// GetUserProfile 获取用户资料
func (s *authService) GetUserProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 清除敏感信息（密码）
	user.Password = ""
	return user, nil
}