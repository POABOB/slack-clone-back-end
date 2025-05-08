package repository

import (
	"github.com/POABOB/slack-clone-back-end/services/user-service/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 創建新的用戶資料存取實例
func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *user.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*user.User, error) {
	var singleUser user.User
	err := r.db.First(&singleUser, id).Error
	if err != nil {
		return nil, err
	}
	return &singleUser, nil
}

func (r *userRepository) FindByEmail(email string) (*user.User, error) {
	var singleUser user.User
	err := r.db.Where("email = ?", email).First(&singleUser).Error
	if err != nil {
		return nil, err
	}
	return &singleUser, nil
}

func (r *userRepository) Update(user *user.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Model(&user.User{}).Where("id = ?", id).Update("is_deleted", true).Error
}
