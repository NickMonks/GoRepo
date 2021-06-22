package services

import (
	"strings"

	"github.com/nickmonks/microservices-go/src/api/config"
	"github.com/nickmonks/microservices-go/src/api/domain/github"
	"github.com/nickmonks/microservices-go/src/api/domain/repositories"
	github_provider "github.com/nickmonks/microservices-go/src/api/providers/github_providers"
	"github.com/nickmonks/microservices-go/src/api/utils/errors"
)

type repoService struct {
}

type repoServiceInterface interface {
	CreateRepo(request repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService repoServiceInterface
)

// called a single time every time is imported
func init() {
	RepositoryService = &repoService{}
}

// This function will create the repo by calling the provider in github
func (s *repoService) CreateRepo(input repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, errors.NewBadRequestError("invalid Repository name")
	}

	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	// the response that we will retrieve to the client
	// (we want to pass this struct)
	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil

}

// create repo concurrently (we called the go routine and create a channel)
func (s *repoService) CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {
	input := make(chan repositories.CreateRepositoriesResult)
	for _, current := range request {
		go s.createRepoConcurrent(current, input)

	}
}

func (s *repoService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{
			Error: err,
		}

		return
	}
	// If no error, create the request to Github API
	request := github.CreateRepoRequest{
		Name:        input.Name,
		Description: input.Description,
		Private:     false,
	}

	result, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		output <- repositories.CreateRepositoriesResult{
			Error: errors.NewApiError(err.StatusCode, err.Message),
		}
	}

	// Once validation is successful, we simply send the repository result to the channel
	output <- repositories.CreateRepositoriesResult{
		Response: &repositories.CreateRepoResponse{
			Id:    result.Id,
			Name:  result.Name,
			Owner: result.Owner.Login,
		},
	}
}
