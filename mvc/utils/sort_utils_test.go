package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBubbleSort(t *testing.T) {
	els := []int{9, 8, 7, 6, 5}

	els = BubbleSort(els)

	assert.NotNil(t, els)
}
