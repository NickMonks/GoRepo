package github_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nickmonks/microservices-go/src/api/client/restclient"
	"github.com/nickmonks/microservices-go/src/api/domain/github"
)

const (
	headerAuthorization       = "Authorization"
	headerAuthorizationFormat = "token %s"

	urlCreateRepo = "https://api.github.com/user/repos"
)

func getAuthorizationHeader(accessToken string) string {
	return fmt.Sprintf(headerAuthorizationFormat, accessToken)
}

func CreateRepo(accessToken string, request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GithubErrorResponse) {

	headers := http.Header{}
	headers.Set(headerAuthorization, getAuthorizationHeader(accessToken))
	response, err := restclient.Post(urlCreateRepo, request, headers)
	if err != nil {
		log.Println(fmt.Sprintf("posting error: %s", err.Error()))
		return nil, &github.GithubErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil && bytes != nil {
		log.Println("Error when trying to create new repo...")
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid response body. Can't read bytes of memory"}
	}

	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(bytes, &errResponse); err != nil {
			return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "invalid json response body"}
		}
		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}

	var result github.CreateRepoResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when unmarshal repo successful response: %s", err.Error()))
		return nil, &github.GithubErrorResponse{StatusCode: http.StatusInternalServerError, Message: "error unmarshalling"}
	}

	return &result, nil
}
