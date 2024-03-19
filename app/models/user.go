package models

import (
	"errors"

	"github.com/Tanakaryuki/brachio-backend/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      string `json:"UserId" gorm:"column:user_id;unique;not null"`
	GithubID    string `json:"GithubId" gorm:"column:github_id;not null;unique"`
	DisplayName string `json:"DisplayName" gorm:"column:display_name"`
	IsPublic    bool   `json:"IsPublic" gorm:"column:is_public;default:true"`
	ImageURL    string `json:"ImageUrl" gorm:"column:image_url"`
}

func CreateUser(user *User) error {
	if err := db.DB.Create(user).Error; err != nil {
		return echo.ErrInternalServerError
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	users := []User{}
	if db.DB.Find(&users).Error != nil {
		return nil, echo.ErrNotFound
	}
	return users, nil
}

func GetUserById(id string) (*User, error) {
	user := User{}
	if err := db.DB.Where("github_id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByUserId(id string) (*User, error) {
	user := User{}
	if err := db.DB.Where("user_id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
