package source

type Source interface {
	TransformToData() Data
}

type DummyFile struct {
	Content string
	Title   string
}

func (d DummyFile) TransformToData() Data {
	return Data{Title: d.Title, Content: d.Content}
}

// intermediate data structure to create Record
type Data struct {
	Title   string   `json:"title"`
	Link    string   `json:"link"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type Id string

// database unit for all sources
type Record struct {
	//unique identifier
	ID Id `json:"id"`
	//title
	Title string `json:"title"`
	//potential link to the source if applicable
	Link string `json:"link"`
	//text content to display on results page
	Content string `json:"content"`
	//map of tokens to their frequency
	TokenFrequency map[string]int `json:"tokenFrequency"`
}

type SearchResult struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Content     string `json:"content"`
	Provider    string `json:"provider"`
	Matches     int    `json:"matches"`
	ContentType string `json:"contentType"`
}
