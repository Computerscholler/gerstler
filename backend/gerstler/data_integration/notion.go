package data_integration

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/kjk/notionapi"

	"github.com/computerscholler/gerstler/source"
)

const notionProvider = "notion"

type NotionResultWrap struct {
	Results   Results
	Total     int
	RecordMap map[string]interface{}
}

type NotionResult struct {
	Id          string
	IsNavigable bool
	Highlight   map[string]string
	Score       float32
}

type Results []NotionResult

func (a Results) Len() int           { return len(a) }
func (a Results) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a Results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type NotionClient struct {
	client  notionapi.CachingClient
	secret  string
	spaceid string
}

func ReadNotionSecret(path string) (string, string, error) {
	f, err := os.Open(path + "notion")
	if err != nil {
		return "", "", fmt.Errorf("notion: unable to find Notion secret: %+v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	token := scanner.Text()
	scanner.Scan()
	spaceid := scanner.Text()
	return token, spaceid, nil
}

func (client NotionClient) Metadata() Metadata {
	return Metadata{Name: notionProvider, DisplayName: "Notion"}
}

func CreateNotionClient(notionSecret string, spaceid string) (DataIntegrator, error) {
	client := notionapi.Client{}
	client.AuthToken = notionSecret

	cacheClient, err := notionapi.NewCachingClient("./CacheDir", &client)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to create caching client %+v", err)
	}
	cacheClient.PreLoadCache()
	return NotionClient{*cacheClient, notionSecret, spaceid}, nil
}

// <gzkNfoUU>VS</gzkNfoUU> Tag is used for highlighting the token
func (client NotionClient) Search(tokens []string) ([]source.SearchResult, error) {

	Data := `{"type":"BlocksInSpace","query":"SEARCH_QUERY","spaceId":"` + client.spaceid + `","limit":8,"filters":{"isDeletedOnly":false,"excludeTemplates":false,"isNavigableOnly":false,"requireEditPermissions":false,"ancestors":[],"createdBy":[],"editedBy":[],"lastEditedTime":{},"createdTime":{}},"sort":"Relevance","source":"quick_find"}`

	Data = strings.Replace(Data, "SEARCH_QUERY", strings.Join(tokens, " "), 1)

	req, err := http.NewRequest("POST", "https://www.notion.so/api/v3/search", strings.NewReader(Data))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cookie", "token_v2="+client.secret)

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Fatalf("Request failed %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	var resultwrap NotionResultWrap
	json.Unmarshal([]byte(sb), &resultwrap)

	sort.Sort(resultwrap.Results)

	var dataResults []source.SearchResult

	for _, result := range resultwrap.Results {
		// TODO get more content from block?
		// block, err := client.client.DownloadPage(result.Id)
		// if err != nil {
		// 	log.Fatalf("Download of Notion page %s failed: %v",result.Id, err)
		// }
		// fmt.Printf("BLOCK:%+v\n", block)
		// TODO remove highlight tags (use for frontend)
		title := removeHighlightTags(result.Highlight["pathText"])
		if title != "" {
			dataResults = append(dataResults, source.SearchResult{Title: title, Link: "block.NotionURL()", Content: removeHighlightTags(result.Highlight["text"]), Provider: notionProvider})
		}

	}

	return dataResults, nil
}

func removeHighlightTags(content string) string {
	content = strings.Replace(content, "<gzkNfoUU>", "", -1)
	content = strings.Replace(content, "</gzkNfoUU>", "", -1)
	return content
}
