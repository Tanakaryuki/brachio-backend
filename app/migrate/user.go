package migrate

import (
	"github.com/Tanakaryuki/brachio-backend/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
}
