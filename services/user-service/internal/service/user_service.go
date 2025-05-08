package service

import (
	"github.com/POABOB/slack-clone-back-end/pkg/auth"
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
)

type userService struct {
	repo user.UserRepository
}

// NewUserService 創建新的使用者服務實例
func NewUserService(repo user.UserRepository) user.UserService {
	return &userService{
		repo: repo,
	}
}

// GetUserByID 獲取使用者訊息
func (s *userService) GetUserByID(id uint) (*user.User, error) {
	return s.repo.FindByID(id)
}

// UpdateUser 更新使用者訊息
func (s *userService) UpdateUser(user *user.User) error {
	// 如果密碼被更新，需要重新加密
	if user.Password != "" {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	return s.repo.Update(user)
}

// DeleteUser 刪除使用者
func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
