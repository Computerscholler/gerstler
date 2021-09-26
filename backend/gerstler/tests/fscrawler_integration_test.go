package main

import (
	"log"
	"testing"

	"github.com/sintemal/gerstler/data_integration"
	"github.com/stretchr/testify/assert"
)

func TestCralwer(t *testing.T) {
	client := data_integration.CreateCrawlerClient()

	output := client.Search([]string{"Hello"})
	log.Printf("%+v\n", output)
	assert.NotEmpty(t, output)
}
