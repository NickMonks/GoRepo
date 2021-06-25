package config

import (
	"fmt"
	"os"
)

const (
	apiGithubAccessToken = "SECRET_GITHUB_ACCESS_TOKEN"
	LogLevel             = "info"
	goEnvironment        = "GO_ENVIRONMENT"
	production           = "production"
)

var (
	githubAccesToken = os.Getenv(apiGithubAccessToken)
)

func GetGithubAccessToken() string {
	fmt.Println(githubAccesToken)
	return githubAccesToken
}

func IsProduction() bool {
	return os.Getenv(goEnvironment) == production
}
