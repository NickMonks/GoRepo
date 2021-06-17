package github_provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("ABC123")
	assert.EqualValues(t, "token ABC123", header)
}
