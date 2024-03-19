package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Tanakaryuki/brachio-backend/models"
	"github.com/Tanakaryuki/brachio-backend/schemas"
)

func GetFollowersByGithubID(githubID string) ([]schemas.Followers, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/followers", githubID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var followers []schemas.Followers
	if err := json.NewDecoder(res.Body).Decode(&followers); err != nil {
		return nil, err
	}
	return followers, nil
}

func InitializeCommit(github_id string) error {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", github_id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var events []schemas.Event
	if err := json.NewDecoder(res.Body).Decode(&events); err != nil {
		return err
	}
	for _, event := range events {
		if event.Type == "PushEvent" {
			for _, commit := range event.Payload.Commits {
				if err := models.CreateEvent(&models.Event{
					UserID: github_id,
					SHA:    commit.SHA,
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
