package github_provider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/nickmonks/microservices-go/src/api/client/restclient"
	"github.com/nickmonks/microservices-go/src/api/domain/github"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Main entry point of testing
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("ABC123")
	assert.EqualValues(t, "token ABC123", header)
}

func TestCreateRepoErrorRestClient(t *testing.T) {
	// Start and Add bespoke mockup
	restclient.FlushMockups()
	restclient.AddMockups(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Err:        errors.New("invalid restclient response"),
	})
	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
}

func TestCreateRepoErrorInvalidResponseBody(t *testing.T) {
	// Start and Add bespoke mockup
	restclient.FlushMockups()

	invalidCloser, _ := os.Open("-asf3")

	restclient.AddMockups(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       invalidCloser,
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid response body. Can't read bytes of memory", err.Message)
}

func TestCreateRepoInvalidErrorInterface(t *testing.T) {
	// Start and Add bespoke mockup
	restclient.FlushMockups()

	restclient.AddMockups(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message":1`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "invalid json response body", err.Message)
}

func TestCreateRepoInvalidSuccessResponse(t *testing.T) {
	// Start and Add bespoke mockup
	restclient.FlushMockups()

	restclient.AddMockups(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":"123"`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.StatusCode)
	assert.EqualValues(t, "error unmarshalling", err.Message)
}

func TestCreateRepoWithNoError(t *testing.T) {
	// Start and Add bespoke mockup
	restclient.FlushMockups()

	restclient.AddMockups(restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id":123, "name": "golang-tutorial","full_name":"nickmonks/microservices-go"}`)),
		},
	})

	response, err := CreateRepo("", github.CreateRepoRequest{})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, 123, response.Id)
	assert.EqualValues(t, "golang-tutorial", response.Name)
	assert.EqualValues(t, "nickmonks/microservices-go", response.FullName)

}
