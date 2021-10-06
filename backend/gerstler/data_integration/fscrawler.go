package data_integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/computerscholler/gerstler/source"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

const NAME = "filesystem"

type CrawlerClient struct {
	client elasticsearch7.Client
}

type SearchResult struct {
	Hits SearchResult2
}
type SearchResult2 struct {
	Hits []Hit
}
type Hit struct {
	Index     string  `json:"_index"`
	Type      string  `json:"_type"`
	Id        string  `json:"_id"`
	Score     float32 `json:"_score"`
	Source    File    `json:"_source"`
	Highlight Highlight
}

type Highlight struct {
	Content []string
}

type File struct {
	File Filename
}
type Filename struct {
	Filename string
}

func CreateCrawlerClient() (DataIntegrator, error) {
	client, err := elasticsearch7.NewClient(elasticsearch7.Config{Addresses: []string{"http://elasticsearch:9200"}})
	if err != nil {
		return nil, fmt.Errorf("fscrawler: cannt create crawler client %+v", err)
	}
	return CrawlerClient{client: *client}, nil
}

func (client CrawlerClient) Search(tokens []string) ([]source.SearchResult, error) {

	searchResults := []source.SearchResult{}

	var queryJson bytes.Buffer

	query := map[string]interface{}{
		"_source": []string{"file.filename"},
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    strings.Join(tokens, " "),
				"fields":   []string{"content", "file.filename"},
				"operator": "or",
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"content": map[string]interface{}{
					"fragment_size": 80,
				},
			},
		},
	}

	if err := json.NewEncoder(&queryJson).Encode(query); err != nil {
		return nil, fmt.Errorf("fscrawler: error encoding query: %s", err)
	}

	// log.Printf(queryJson.String())

	res, err := client.client.Search(
		client.client.Search.WithContext(context.Background()),
		client.client.Search.WithBody(&queryJson),
		client.client.Search.WithTrackTotalHits(true),
		client.client.Search.WithPretty(),
	)

	if err != nil {
		return nil, fmt.Errorf("fscrawler: error while searching %+v", err)
	}

	respBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("fscrawler: couldn't read body %+v", err)
	}

	var hits SearchResult

	json.Unmarshal(respBody, &hits)

	for _, hit := range hits.Hits.Hits {
		searchResults = append(searchResults, source.SearchResult{Title: hit.Source.File.Filename, Link: "", Content: hit.Highlight.Content[0], Provider: NAME})
	}
	return searchResults, nil
}

func (client CrawlerClient) Metadata() Metadata {
	return Metadata{Name: NAME, DisplayName: "FileSystem"}
}
