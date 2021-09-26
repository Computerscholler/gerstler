package data_integration_test

import (
	"fmt"
	"testing"

	"github.com/sintemal/gerstler/data_integration"
)
func TestNotion(t *testing.T){
	secret,spaceId := data_integration.ReadNotionSecret("../../../secrets/")
	notionClient := data_integration.CreateNotionClient(secret,spaceId)
	res := notionClient.Search([]string{"vs","code"})
	fmt.Println(res)
}	
