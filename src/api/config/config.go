package config

import (
	"fmt"
	"os"
)

const (
	apiGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
)

var (
	githubAccesToken = os.Getenv(apiGithubAccessToken)
)

func GetGithubAccessToken() string {
	fmt.Println(githubAccesToken)
	return githubAccesToken
}
