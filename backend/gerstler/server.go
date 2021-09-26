package gerstler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sintemal/gerstler/data_integration"
	fs "github.com/sintemal/gerstler/filesystem"
	"github.com/sintemal/gerstler/search"
	"github.com/sintemal/gerstler/source"
)

const (
	ENV_SECRETS_PATH = "GERSTLER_SECRETS_PATH"
	ENV_PORT         = "GERSTLER_PORT"
)

const (
	ACTION_RESULTS        = "results"
	ACTION_LOADING_STATUS = "loading_status"
)

type Server struct {
	SearchClients []data_integration.DataIntegrator
}

var indexPath = "./data/index.json"
var recordsPath = "./data/records.json"
var secretPath = "../../secrets/"

const ActionQuery = "query"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func initInverseIndexer() {
	recFile, _ := fs.OpenFile(recordsPath)
	indexFile, _ := fs.OpenFile(indexPath)
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

func getEnv(name string, def string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		return def
	}

	return val
}

func Start() {
	port, err := strconv.Atoi(getEnv(ENV_PORT, "5000"))
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s has invalid value: %s", ENV_PORT, err))
	}

	initInverseIndexer()

	if path := os.Getenv(ENV_SECRETS_PATH); len(path) > 0 {
		secretPath = path
	}

	gdriveClient := data_integration.CreateClient(secretPath)
	secret, spaceId := data_integration.ReadNotionSecret(secretPath)
	notionClient := data_integration.CreateNotionClient(secret, spaceId)
	fsCrawler := data_integration.CreateCrawlerClient()
	// change credentials
	emailClient, err := data_integration.CreateEmailClient(data_integration.EmailConfig{
		URL:      "imap.example.com",
		Port:     993,
		Username: "john@example.com",
		Password: "password",
	})
	if err != nil {
		log.Fatalln("Failed to create email client:", err)
	}

	searchClients := []data_integration.DataIntegrator{gdriveClient, fsCrawler, notionClient, emailClient}

	server := Server{
		SearchClients: searchClients,
	}

	router := gin.Default()
	router.GET("/api/ws", server.wsHandler)
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router)
	if err != nil {
		log.Fatalln("Server encountered an error:", err)
	}
}
