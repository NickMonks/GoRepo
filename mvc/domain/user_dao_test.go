package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// We need to implement the interface testing.T
// in go, is very important NOT TO test the same function with the other return types
// you should return only for one of the returns.
func TestGetUserNoUserFound(t *testing.T) {
	// Given, or Initialization
	// When, or Execution
	// Then, or Validation
	user, err := GetUser(0)

	assert.Nil(t, user, "we were not expecting a user with id 0")
	assert.NotNil(t, err, "we were expecting an error")
	assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	assert.EqualValues(t, "not_found", err.Code)
	assert.EqualValues(t, "user 0 was not found", err.Message)

}

func TestGetUserNoError(t *testing.T) {
	user, err := GetUser(123)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 123, user.Id)
	assert.EqualValues(t, "Fede", user.FirstName)

}
