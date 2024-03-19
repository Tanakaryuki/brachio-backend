package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Tanakaryuki/brachio-backend/models"
	"github.com/Tanakaryuki/brachio-backend/schemas"
)

var ExtensionToLanguage = map[string]string{
	".as":     "AngelScript",
	".s":      "Assembly",
	".bal":    "Ballerina",
	".c":      "C",
	".h":      "C",
	".cc":     "C++",
	".cp":     "C++",
	".cpp":    "C++",
	".cxx":    "C++",
	".cobol":  "COBOL",
	".coffee": "CoffeeScript",
	".css":    "CSS",
	".clj":    "Clojure",
	".d":      "Makefile",
	".dart":   "Dart",
	".erl":    "Erlang",
	".forth":  "Forth",
	".f":      "Forth",
	".for":    "Forth",
	".f95":    "Fortran",
	".f90":    "Fortran",
	".f03":    "Fortran",
	".go":     "Go",
	".groovy": "Groovy",
	".hs":     "Haskell",
	".lhs":    "Haskell",
	".hx":     "Haxe",
	".html":   "HTML",
	".htm":    "HTML",
	".xhtml":  "XHTML",
	".java":   "Java",
	".js":     "JavaScript",
	".jsx":    "JavaScript",
	".kt":     "Kotlin",
	".lisp":   "CommonLisp",
	".lsl":    "LSL",
	".lua":    "Lua",
	".mat":    "MATLAB",
	".m":      "Objective-C",
	".ml":     "OCaml",
	".pas":    "Pascal",
	".pl":     "Perl",
	".php":    "PHP",
	".pro":    "Prolog",
	".py":     "Python",
	".r":      "R",
	".rb":     "Ruby",
	".rs":     "Rust",
	".scala":  "Scala",
	".scm":    "Scheme",
	".sh":     "Shell",
	".cs":     "Smalltalk",
	".swift":  "Swift",
	".ts":     "TypeScript",
	".vbs":    "Visual Basic",
	".xquery": "XQuery",
	".v":      "Coq",
}

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

func UpdateCommit(github_id string) error {
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
				if event, err := models.GetEventBySHA(commit.SHA); err == nil {
					if event == nil {
						if err := GetCommitDetailByURL(commit.URL, commit.Author.Name); err != nil {
							return err
						}
						err = models.CreateEvent(&models.Event{
							UserID: github_id,
							SHA:    commit.SHA,
						})
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func GetCommitDetailByURL(url string, UserID string) error {
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

	var file schemas.CommitInfo
	if err := json.NewDecoder(res.Body).Decode(&file); err != nil {
		return err
	}
	for _, file := range file.Files {
		parts := strings.Split(file.Filename, ".")
		lastPart := parts[len(parts)-1]
		lastPart = "." + lastPart
		Language := ExtensionToLanguage[lastPart]

		if Language != "" {
			if err := models.UpDatePet(UserID, Language, file.Changes); err != nil {
				return err
			}
		}
	}

	return nil
}
