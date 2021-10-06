package data_integration

import (
	"github.com/computerscholler/gerstler/source"
)

type DataIntegrator interface {
	Search([]string) ([]source.SearchResult, error)
	Metadata() Metadata
}

type Metadata struct {
	Name        string
	DisplayName string
}
