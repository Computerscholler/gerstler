package gerstler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/computerscholler/gerstler/data_integration"
	fs "github.com/computerscholler/gerstler/filesystem"
	"github.com/computerscholler/gerstler/search"
	"github.com/computerscholler/gerstler/source"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	SECRETS_PATH = "secrets_path"
	PORT         = "server.port"
	INDEX_PATH   = "index_path"
	RECORDS_PATH = "records_path"
	EMAIL        = "email"
)

const (
	ACTION_RESULTS        = "results"
	ACTION_LOADING_STATUS = "loading_status"
)

type Server struct {
	SearchClients []data_integration.DataIntegrator
}

const ActionQuery = "query"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func initInverseIndexer() {
	recFile, _ := fs.OpenFile(viper.GetString(RECORDS_PATH))
	indexFile, _ := fs.OpenFile(viper.GetString(INDEX_PATH))
	source.LoadDb(recFile, indexFile)
	source.GenerateIndexer()
}

type PayloadRequest struct {
	Query  string `json:"query"`
	Action string `json:"action"`
}

type LoadingStatus struct {
	Provider string `json:"provider"`
	Loading  bool   `json:"loading"`
}

type PayloadResponse struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
	Query  string      `json:"query"`
}

func (s *Server) sendSearchResults(conn *websocket.Conn, query string, writer chan PayloadResponse) {
	queries := search.Analyze(query)
	for _, client := range s.SearchClients {
		go func(dataclient data_integration.DataIntegrator) {
			writer <- PayloadResponse{Action: ACTION_LOADING_STATUS, Data: LoadingStatus{Loading: true, Provider: dataclient.Metadata().Name}, Query: query}
			res, err := dataclient.Search(queries)
			if err != nil {
				log.Printf("Error: %s: %s\n", dataclient.Metadata().Name, err)
			} else {
				results_resp := PayloadResponse{Action: ACTION_RESULTS, Data: res, Query: query}
				writer <- results_resp
			}
			writer <- PayloadResponse{Action: ACTION_LOADING_STATUS, Data: LoadingStatus{Loading: false, Provider: dataclient.Metadata().Name}, Query: query}

		}(client)
	}
}

func (s *Server) readRequest(conn *websocket.Conn) (*PayloadRequest, error) {
	req := &PayloadRequest{}
	err := conn.ReadJSON(req)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Fatalf("unexpected close error: %v\n", err)
		}
		return nil, err
	}
	return req, nil
}

// handle search results
func (s *Server) wsHandler(ctx *gin.Context) {
	// POST quick search
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	writer := make(chan PayloadResponse, 5)
	done := make(chan bool, 1)

	go func() {
		for {
			select {
			case resp := <-writer:
				err = conn.WriteJSON(&resp)
				if err != nil {
					log.Printf("Could not write JSON (%+v), result: %v", resp, err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		req, err := s.readRequest(conn)
		if err != nil {
			break
		}
		log.Printf("Request: %s\n", req.Query)
		s.sendSearchResults(conn, req.Query, writer)
	}
	done <- true
}

func loadConfig() {
	viper.SetConfigName("gerstler")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("/app/config")
	viper.SetEnvPrefix("GERSTLER")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// set defaults
	viper.SetDefault(PORT, "5000")
	viper.SetDefault(SECRETS_PATH, "../../secrets/")
	viper.SetDefault(INDEX_PATH, "./data/index.json")
	viper.SetDefault(RECORDS_PATH, "./data/records.json")
}

func Start() {
	loadConfig()
	initInverseIndexer()

	secretPath := viper.GetString(SECRETS_PATH)

	searchClients := make([]data_integration.DataIntegrator, 0, 5)

	gdriveClient, err := data_integration.CreateClient(secretPath)
	if err == nil {
		searchClients = append(searchClients, gdriveClient)
	} else {
		log.Println(err)
	}

	secret, spaceId, err := data_integration.ReadNotionSecret(secretPath)
	if err == nil {
		notionClient, err := data_integration.CreateNotionClient(secret, spaceId)
		if err == nil {
			searchClients = append(searchClients, notionClient)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	fsCrawler, err := data_integration.CreateCrawlerClient()
	if err == nil {
		searchClients = append(searchClients, fsCrawler)
	} else {
		log.Println(err)
	}

	emailConfig := viper.Sub(EMAIL)
	if emailConfig != nil {
		emailClient, err := data_integration.CreateEmailClientFromViper(emailConfig)
		if err == nil {
			searchClients = append(searchClients, emailClient)
		} else {
			log.Println("failed to create email client:", err)
		}
	}

	server := Server{
		SearchClients: searchClients,
	}

	router := gin.Default()
	router.GET("/api/ws", server.wsHandler)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", viper.GetInt(PORT)), router)
	if err != nil {
		log.Fatalln("Server encountered an error:", err)
	}
}
