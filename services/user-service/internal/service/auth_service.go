package service

import (
	"errors"

	authlib "github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt"
	"github.com/POABOB/slack-clone-back-end/pkg/auth/jwt/rbac"

	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/auth"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
)

type authService struct {
	userRepo   user.UserRepository
	jwtManager jwt.TokenManager
}

// NewAuthService 創建新的驗證服務實例
func NewAuthService(userRepo user.UserRepository, jwtManager jwt.TokenManager) auth.AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register 註冊新使用者
func (s *authService) Register(user *user.User) error {
	// 檢查 Email 是否存在
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	// 加密密碼
	hashedPassword, err := authlib.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.Create(user)
}

// Login 使用者登入
func (s *authService) Login(email, password string) (string, error) {
	// 查找使用者
	singleUser, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 驗證密碼
	err = authlib.CheckPassword(password, singleUser.Password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	tokenString, err := s.GenerateToken(singleUser)
	return tokenString, err
}

// GenerateToken 產生 JWT Token
func (s *authService) GenerateToken(singleUser *user.User) (string, error) {
	token, err := s.jwtManager.GenerateToken(rbac.NewRBACClaims(
		singleUser.ID,
		singleUser.Email,
		singleUser.Username,
		singleUser.Role,
		singleUser.Permissions,
	))
	if err != nil {
		return "", err
	}
	return token, nil
}
