package routes

import (
	"net/http"

	"github.com/Tanakaryuki/brachio-backend/firebase"
	"github.com/Tanakaryuki/brachio-backend/models"
	"github.com/labstack/echo/v4"
)

func createUser(context echo.Context) error {
	userID, err := firebase.VerifyTokenAndGetUserID(context.Request())
	if err != nil {
		return context.JSON(http.StatusUnauthorized, err.Error())
	}
	user := new(models.User)
	user.UserID = userID
	if err := context.Bind(user); err != nil {
		return context.JSON(http.StatusBadRequest, "ユーザーを作成できませんでした。")
	}

	if err := models.CreateUser(user); err != nil {
		return context.JSON(http.StatusInternalServerError, "ユーザーを作成できませんでした。")
	}
	return context.JSON(http.StatusCreated, user)
}

func getAllUsers(context echo.Context) error {
	users, err := models.GetAllUsers()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "ユーザーを取得できませんでした。")
	}
	return context.JSON(http.StatusOK, users)
}

func getUserById(context echo.Context) error {
	userId := context.Param("userId")
	user, err := models.GetUserById(userId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "ユーザーを取得できませんでした。")
	} else if user == nil {
		return context.JSON(http.StatusNotFound, "ユーザーが見つかりませんでした。")
	}
	return context.JSON(http.StatusOK, user)
}
