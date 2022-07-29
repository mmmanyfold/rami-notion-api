package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomei/notionapi"
)

var database = map[string]notionapi.DatabaseID{
	"projects":   notionapi.DatabaseID("bee593efdc654282911f3dc5550e144a"),
	"homepage":   notionapi.DatabaseID("a79aece399014bc282a27024de23464a"),
	"transcripts": notionapi.DatabaseID("d815aa37777a4b04812f38b0b9d81b89"),
}

func main() {
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}
	db, err := client.Database.Query(context.Background(), database["projects"], &dbRequest)
	if err != nil {
		panic(err)
	}

	if len(db.Results) > 0 {
		var pages []notionapi.Page
		for i, r := range db.Results {
			fmt.Printf("id: %s\n", r.ID)
			pages = append(pages, db.Results[i])
		}
	}
}
