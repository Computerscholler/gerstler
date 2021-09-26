package main

import (
	"log"

	"github.com/sintemal/gerstler/data_integration"
)

func main() {
	searchTerm := []string{"Hello World"}
	client := data_integration.CreateClient("../../secrets/")
	log.Printf("%v", client.Search(searchTerm))

}
