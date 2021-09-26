package data_integration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sintemal/gerstler/source"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const gdriveProvider = "gdrive"

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, path string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := filepath.Join(path, "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	reader := bufio.NewReader(os.Stdin)

	authCode, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

type GoogleDriveClient struct {
	drive       *drive.Service
	credentials []byte
	metatdata   Metadata
}

func CreateClient(path string) DataIntegrator {
	var gdriveCredentialsPath = filepath.Join(path, "gdrive.json")

	credentials, err := ioutil.ReadFile(gdriveCredentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	ctx := context.Background()
	config, err := google.ConfigFromJSON(credentials, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, path)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return GoogleDriveClient{srv, credentials, Metadata{DisplayName: "Google Drive Client", Name: gdriveProvider}}
}

func (client GoogleDriveClient) Search(tokens []string) ([]source.SearchResult, error) {
	query := "fullText contains '" + strings.Join(tokens, "' or fullText contains '") + "'"
	r, err := client.drive.Files.List().Q(query).Fields("files(*)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve search Result: %v", err)
	}
	return client.parseFiles(r.Files, tokens), nil
}

func (client GoogleDriveClient) parseFiles(files []*drive.File, tokens []string) []source.SearchResult {

	datas := []source.SearchResult{}
	for _, file := range files {
		// var content *http.Response
		// var err error
		// file.Content
		data := source.SearchResult{Title: file.Name, Link: file.WebViewLink, Content: "", Provider: gdriveProvider}
		datas = append(datas, data)

		// var document []byte
		// if err != nil {
		// 	log.Fatalf("Unable to get file %v", err)
		// }

		// content.Location()
		// log.Printf("%+v\n", content)
		// log.Printf("Data len: %v", content.ContentLength)

		// _, err = content.Body.Read(document)

		// if err != nil {s
		// 	log.Fatalf("Unable to read file %v", err)
		// }
		// data := source.SearchResult{Title: file.Name, Link: file.WebViewLink, Content: string(document), Provider: providerType}
	}

	return datas
}

func (client GoogleDriveClient) Metadata() Metadata {
	return client.metatdata
}
