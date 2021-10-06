package main

import (
	"testing"

	"github.com/computerscholler/gerstler"
	"github.com/computerscholler/gerstler/source"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

type Payload struct {
	Query  string `json:"query"`
	Action string `json:"action"`
}

func TestWsConnection(t *testing.T) {
	gerstler.Start()
	wsURL := "ws" + "://localhost:5000" + "/api/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	defer ws.Close()
	req := Payload{"hi", "query"}
	ws.WriteJSON(&req)
	var got []source.SearchResult
	err = ws.ReadJSON(&got)
	assert.NoError(t, err)
	assert.Equal(t, []source.SearchResult{}, got)
}
