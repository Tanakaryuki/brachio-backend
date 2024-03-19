package routes

import (
	"net/http"

	"github.com/Tanakaryuki/brachio-backend/firebase"
	"github.com/Tanakaryuki/brachio-backend/github"
	"github.com/Tanakaryuki/brachio-backend/models"
	"github.com/Tanakaryuki/brachio-backend/schemas"
	"github.com/labstack/echo/v4"
)

func getPetsByUserId(context echo.Context) error {
	userId := context.Param("userId")
	user, _ := models.GetUserById(userId)
	if user == nil {
		return context.JSON(http.StatusNotFound, "ユーザーが見つかりませんでした。")
	}
	if err := github.UpdateCommit(userId); err != nil {
		return context.JSON(http.StatusInternalServerError, "コミット情報を取得できませんでした。")
	}

	me := schemas.User{
		GithubID:    user.GithubID,
		DisplayName: user.DisplayName,
		ImageURL:    user.ImageURL,
	}
	pets, err := models.GetPetsByUserId(userId)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "ペットを取得できませんでした。")
	}
	followers, err := github.GetFollowersByGithubID(user.GithubID)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "フォロワーを取得できませんでした。(笑)")
	}

	hoge := map[string]interface{}{
		"user":      me,
		"pets":      pets,
		"followers": followers,
	}
	return context.JSON(http.StatusOK, hoge)
}

func feedPet(context echo.Context) error {
	userID, err := firebase.VerifyTokenAndGetUserID(context.Request())
	if err != nil {
		return context.JSON(http.StatusUnauthorized, err.Error())
	}
	user, _ := models.GetUserByUserId(userID)
	if user == nil {
		return context.JSON(http.StatusNotFound, "ユーザーが見つかりませんでした。")
	}
	language := context.Param("language")
	pet, err := models.GetPetByLanguage(user.GithubID, language)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "ペットを餌やりできませんでした。")
	}
	if err := models.FeedPet(pet); err != nil {
		return context.JSON(http.StatusInternalServerError, "ペットを餌やりできませんでした。")
	}

	if err := github.UpdateCommit(user.GithubID); err != nil {
		return context.JSON(http.StatusInternalServerError, "コミット情報を取得できませんでした。")
	}

	me := schemas.User{
		GithubID:    user.GithubID,
		DisplayName: user.DisplayName,
		ImageURL:    user.ImageURL,
	}
	pets, err := models.GetPetsByUserId(user.GithubID)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "ペットを取得できませんでした。")
	}
	followers, err := github.GetFollowersByGithubID(user.GithubID)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, "フォロワーを取得できませんでした。(笑)")
	}

	hoge := map[string]interface{}{
		"user":      me,
		"pets":      pets,
		"followers": followers,
	}
	return context.JSON(http.StatusOK, hoge)
}
