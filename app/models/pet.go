package models

import (
	"errors"
	"time"

	"github.com/Tanakaryuki/brachio-backend/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Event struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID string `json:"user_id"`
	SHA    string `json:"sha"`
}

type Pet struct {
	ID              uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          string `json:"user_id"`
	Language        string `json:"Language"`
	HungerLevel     int    `json:"HungerLevel"`
	FriendshipLevel int    `json:"FriendshipLevel"`
	EscapeNum       int    `json:"EscapeNum"`
	BaitsNum        int    `json:"BaitsNum"`
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

func GetPetsByUserId(userId string) ([]*Pet, error) {
	pets := make([]*Pet, 0)
	if err := db.DB.Where("user_id = ?", userId).Find(&pets).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return pets, nil
}

func GetPetByLanguage(userId string, language string) (*Pet, error) {
	pet := Pet{}
	if err := db.DB.Where("user_id = ? AND language = ?", userId, language).First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pet, nil
}

func UpDatePet(UserID string, Language string, FriendshipLevel int) error {
	hoge, err := GetPetByLanguage(UserID, Language)
	if err != nil {
		return err
	}

	if hoge == nil {
		pet := &Pet{
			UserID:          UserID,
			Language:        Language,
			HungerLevel:     100,
			FriendshipLevel: FriendshipLevel,
			EscapeNum:       0,
			BaitsNum:        0,
		}
		if err := db.DB.Create(pet).Error; err != nil {
			return err
		}

		go func() {
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()

			for {
				select {
				case <-timer.C:
					pet, err := DecreaseHungerLevel(UserID, Language)
					if err != nil {
						continue
					}
					if pet.HungerLevel == 0 {
						ResetPet(UserID, Language)
						return
					}
					timer.Reset(10 * time.Second)
				}
			}
		}()
	} else if hoge.HungerLevel == 0 {
		pet := &Pet{
			UserID:          UserID,
			Language:        Language,
			HungerLevel:     100,
			FriendshipLevel: FriendshipLevel,
			EscapeNum:       hoge.EscapeNum,
			BaitsNum:        0,
		}
		if err := db.DB.Model(&pet).Where("user_id = ? AND language = ?", UserID, Language).Select("*").Updates(pet).Error; err != nil {
			return err
		}

		go func() {
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()

			for {
				select {
				case <-timer.C:
					pet, err := DecreaseHungerLevel(UserID, Language)
					if err != nil {
						continue
					}
					if pet.HungerLevel == 0 {
						ResetPet(UserID, Language)
						return
					}
					timer.Reset(10 * time.Minute)
				}
			}
		}()
	} else {
		pet := &Pet{
			UserID:          UserID,
			Language:        Language,
			HungerLevel:     0,
			FriendshipLevel: FriendshipLevel + hoge.FriendshipLevel,
			EscapeNum:       hoge.EscapeNum,
			BaitsNum:        0,
		}
		if hoge.HungerLevel == 100 {
			pet.HungerLevel = 100
			pet.BaitsNum = hoge.BaitsNum + 1
		} else {
			pet.HungerLevel = hoge.HungerLevel + 10
			if pet.HungerLevel > 100 {
				pet.HungerLevel = 100
			}
		}
		if err := db.DB.Model(&pet).Where("user_id = ? AND language = ?", UserID, Language).Updates(pet).Error; err != nil {
			return err
		}
	}
	return nil
}

func FeedPet(pet *Pet) error {
	if pet.HungerLevel != 100 {
		pet.HungerLevel += 10
		if pet.HungerLevel > 100 {
			pet.HungerLevel = 100
		}
		pet.BaitsNum -= 1
		if err := db.DB.Model(&pet).Where("id = ?", pet.ID).Updates(pet).Error; err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("ペットはお腹いっぱいです。")
	}
}

func DecreaseHungerLevel(UserID string, Language string) (*Pet, error) {
	pet, err := GetPetByLanguage(UserID, Language)
	if err != nil {
		return nil, err
	}
	pet.HungerLevel -= 1
	if err := db.DB.Model(&pet).Where("user_id = ? AND language = ?", UserID, Language).Updates(pet).Error; err != nil {
		return nil, err
	}
	return pet, nil
}

func ResetPet(UserID string, Language string) error {
	pet, err := GetPetByLanguage(UserID, Language)
	if err != nil {
		return err
	}
	pet.HungerLevel = 0
	pet.FriendshipLevel = 0
	pet.BaitsNum = 0
	pet.EscapeNum += 1
	if err := db.DB.Model(&pet).Where("user_id = ? AND language = ?", UserID, Language).Select("*").Updates(pet).Error; err != nil {
		return err
	}
	return nil
}
