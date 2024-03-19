package routes

import (
	"net/http"

	"github.com/Tanakaryuki/brachio-backend/firebase"
	"github.com/Tanakaryuki/brachio-backend/github"
	"github.com/Tanakaryuki/brachio-backend/models"
	"github.com/Tanakaryuki/brachio-backend/schemas"
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

	me := schemas.User{
		GithubID:    user.GithubID,
		DisplayName: user.DisplayName,
		ImageURL:    user.ImageURL,
	}
	if err := github.InitializeCommit(me.GithubID); err != nil {
		return context.JSON(http.StatusInternalServerError, "ユーザー情報を取得できませんでした。")
	}
	pets := make([]schemas.Pet, 0)
	followers, err := github.GetFollowersByGithubID(user.GithubID)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "フォロワーを取得できませんでした。(笑)")
	}

	hoge := map[string]interface{}{
		"user":      me,
		"pets":      pets,
		"followers": followers,
	}
	return context.JSON(http.StatusCreated, hoge)
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
