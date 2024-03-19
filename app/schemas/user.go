package schemas

type User struct {
	GithubID    string `json:"github_id"`
	DisplayName string `json:"display_name"`
	ImageURL    string `json:"avatar_url"`
}

type Followers struct {
	GithubID string `json:"login"`
	ImageURL string `json:"avatar_url"`
}
