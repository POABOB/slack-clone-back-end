package user

import (
	"time"
)

// User 使用者實體
type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	Password    string    `json:"-" gorm:"not null"`
	Username    string    `json:"username" gorm:"not null"`
	Role        string    `json:"role" gorm:"not null;default:'user'"`
	Permissions []string  `json:"permissions" gorm:"type:json"`
	LastLogin   time.Time `json:"last_login"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRepository 使用者資料存取介面
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserService 使用者業務邏輯介面
type UserService interface {
	GetUserByID(id uint) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
}
