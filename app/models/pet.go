package models

import (
	"errors"

	"github.com/Tanakaryuki/brachio-backend/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Event struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID string `json:"user_id"`
	SHA    string `json:"sha"`
}

func CreateEvent(event *Event) error {
	if err := db.DB.Create(event).Error; err != nil {
		return echo.ErrInternalServerError
	}
	return nil
}

func GetEventBySHA(sha string) (*Event, error) {
	event := Event{}
	if err := db.DB.Where("sha = ?", sha).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}
