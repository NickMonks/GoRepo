package services

import (
	"net/http"
	"strings"
	"sync"

	"github.com/nickmonks/microservices-go/src/api/config"
	"github.com/nickmonks/microservices-go/src/api/domain/github"
	"github.com/nickmonks/microservices-go/src/api/domain/repositories"
	"github.com/nickmonks/microservices-go/src/api/log"
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

	log.Info("Sending request to external Api...", "client_id: id", "status: pending")
	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {
		log.Error("Sending request to external Api...", err, "client_id: id", "status: Error")
		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}
	log.Info("Sending request to external Api...", "client_id: id", "status: success")

	// the response that we will retrieve to the client
	// (we want to pass this struct)
	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil

}

// create repo concurrently (we called the go routine and create a channel), TX
func (s *repoService) CreateRepos(request []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {

	// As soon as we get a request from the controller, we create two channels: input and output.
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)

	// defer will close any channels open (in this case the output; input will be closed previosly)
	defer close(output)

	// WaitGroup blocks the until some work has been done.
	var wg sync.WaitGroup

	// Listens to requests
	go s.handleRepoResults(&wg, input, output)

	for _, current := range request {
		wg.Add(1)
		go s.createRepoConcurrent(current, input)

	}

	// the waitgroup will pause until all goroutines have been finished (we add 1 for each call to goroutine)
	// we do this with wg.Done(), written inside handleRepoResults. Once the requests have been done (wg is 0), we will close
	// input channel, terminating the for loop of handleRepoResults
	wg.Wait()
	close(input)

	// once wg has been released, will listen to output channel, which is handleRepoResults, all repos synchronized.
	// because channel input is closed, we exit the for loop and sent the output results, which are appended
	result := <-output

	// Perform some validations for repo creations
	successCreations := 0
	for _, current := range result.Results {
		if current.Response != nil {
			successCreations++
		}
	}

	if successCreations == 0 {
		result.StatusCode = result.Results[0].Error.Status()
	}

	if successCreations == len(request) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}

	return result, nil
}

// This function will produce github requests
func (s *repoService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	// Here, we will listen and receive any incoming events from the CreateRepos, so it will wait
	// and run createRepoConcurrent, producing the input
	for incomingEvent := range input {
		repoResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}

		// finally, we append the results, which is an array of repositories
		results.Results = append(results.Results, repoResult)

		// release the wg by one; so every time createRepoConcurrent sends something to output
		wg.Done()
		// once it reaches the limit (e.g., 3 request), it will go back to the wg.wait() and continue execution
	}

	// input channel is closed; we send all results to the output channel
	output <- results
}

// This function will produce the inputs to be consumed
func (s *repoService) createRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {
	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{
			Error: err,
		}

		return
	}

	result, err := s.CreateRepo(input)

	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}

	// Once validation is successful, we simply send the repository result to the channel
	output <- repositories.CreateRepositoriesResult{
		Response: result,
	}
}
