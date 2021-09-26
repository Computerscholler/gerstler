package search_test

import (
	"testing"

	"github.com/sintemal/gerstler/search"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	data := "Hi world, I am Adrian and he wants to learn Golang."
	tokens := search.Analyze(data)
	assert.Equal(t,[]string{"hi","world","adrian","want","learn","golang"},tokens)
}

