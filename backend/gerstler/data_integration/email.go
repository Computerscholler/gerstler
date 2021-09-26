package data_integration

//Imports
import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/sintemal/gerstler/source"
)

const emailProvider = "email"

type EmailConfig struct {
	URL      string
	SSL      bool
	Port     int
	Username string
	Password string
}

type EmailClient struct {
	Config EmailConfig
	client *client.Client
}

func CreateEmailClient(config EmailConfig) (DataIntegrator, error) {

	client, err := connectEmail(&config)
	if err != nil {
		return nil, err
	}

	return EmailClient{Config: config, client: client}, nil
}

func (ec EmailClient) Metadata() Metadata {
	return Metadata{Name: emailProvider, DisplayName: "E-Mail"}
}

func connectEmail(config *EmailConfig) (*client.Client, error) {
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", config.URL, config.Port), nil)
	if err != nil {
		return nil, err
	}

	if err := c.Login(config.Username, config.Password); err != nil {
		return nil, err
	}

	return c, nil
}

type part struct {
}

func (ec EmailClient) Search(tokens []string) ([]source.SearchResult, error) {
	if ec.client.State() == imap.NotAuthenticatedState || ec.client.State() == imap.LogoutState {
		c, err := connectEmail(&ec.Config)
		if err != nil {
			return nil, err
		}
		ec.client = c
	}

	// Nachrichten aus Inbox laden
	_, err := ec.client.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	// Kriterien setzen
	criteria := imap.NewSearchCriteria()
	criteria.Text = tokens
	ids, err := ec.client.Search(criteria)
	if err != nil {
		return nil, err
	}

	res := []source.SearchResult{}
	// Gefundene Ids verarbeiten
	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		messages := make(chan *imap.Message)
		done := make(chan error, 1)
		go func() {
			done <- ec.client.Fetch(seqset, []imap.FetchItem{imap.FetchBody, imap.FetchBodyStructure, imap.FetchEnvelope, imap.FetchRFC822, imap.FetchRFC822Header, imap.FetchRFC822Size, imap.FetchRFC822Text}, messages)
		}()

		for msg := range messages {

			// Jede message verarbeiten für die Rückgabe
			buf := ""
			var section imap.BodySectionName
			section.Peek = true
			r := msg.GetBody(&section)
			if r == nil {
				log.Println("Server didn't returned message body")
				continue
			}
			// Create a new mail reader
			mr, err := mail.CreateReader(r)
			if err != nil {
				return nil, err
			}

			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, err
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message's text (can be plain-text or HTML)
					b, _ := ioutil.ReadAll(p.Body)
					buf = buf + string(b)
				}
			}

			// Hinzufügen der Message zum Ergebnisarray
			res = append(res, source.SearchResult{
				Title:       msg.Envelope.Subject,
				Link:        "",
				Content:     buf,
				ContentType: "html",
				Provider:    emailProvider,
			})
		}
	}

	return res, nil
}
