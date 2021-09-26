package main

import (
	"testing"

	"github.com/sintemal/gerstler/source"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	data := &source.DummyFile{Title: "Hello", Content: "Hola mundo"}
	data2 := &source.DummyFile{Title: "Dict", Content: "Mundo means world in spanish."}
	source.AddEntry(data)
	source.AddEntry(data2)
	source.GenerateIndexer()
	records, _ := source.Search("mundo")
	assert.Equal(t, 2, len(records))
}
