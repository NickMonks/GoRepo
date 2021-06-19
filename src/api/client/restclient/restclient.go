package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	enableMocks = false
	// set response as a pointer to return nil
	mocks = make(map[string]*Mock)
)

// --- MOCK DEFINITION ----
type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func GetMockId(httpMethod string, urls string) string {
	return fmt.Sprintf("%s_%s", httpMethod, urls)
}

func StartMockups() {
	enableMocks = true
}

func StopMockups() {
	enableMocks = false
}

func AddMockups(mock Mock) {
	mocks[GetMockId(mock.HttpMethod, mock.Url)] = &mock
}

func FlushMockups() {
	// creates a clean slate of mock
	mocks = make(map[string]*Mock)
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {

	if enableMocks {
		mock := mocks[GetMockId(http.MethodPost, url)]
		if mock == nil {
			return nil, errors.New("No mockup found for given resquest")
		}
		return mock.Response, mock.Err
		//TODO: return local mock without calling ext resources
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers

	client := http.Client{}
	return client.Do(request)
}
