package data_integration

import (
	"github.com/sintemal/gerstler/source"
)

type DataIntegrator interface {
	Search([]string) ([]source.SearchResult, error)
	Metadata() Metadata
}

type Metadata struct {
	Name        string
	DisplayName string
}
