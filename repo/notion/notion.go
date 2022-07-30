package notion

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
	"time"
)

// request per second to notion api limit
const rate uint64 = 3

var database = map[string]notionapi.DatabaseID{
	"projects":    notionapi.DatabaseID("bee593efdc654282911f3dc5550e144a"),
	"homepage":    notionapi.DatabaseID("a79aece399014bc282a27024de23464a"),
	"transcripts": notionapi.DatabaseID("d815aa37777a4b04812f38b0b9d81b89"),
}

func GetHomePageProjectsAndAssets(client *notionapi.Client) error {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	ctx := context.Background()
	rateLimiter, err := NewRateLimiter(ctx, "projects", rate, false)
	if err != nil {
		return err
	}

	_, _ = rateLimiter.Take()
	_, _ = rateLimiter.Take()
	_, _ = rateLimiter.Take()
	_, _ = rateLimiter.Take()

	fmt.Println("sleeping for 1 sec ....")
	time.Sleep(time.Second)

	ok, err := rateLimiter.Take()
	if err != nil {
		return err
	}

	if ok {
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

	var _ []rami.HomePageAsset

	defer rateLimiter.Close()
	return nil
}
