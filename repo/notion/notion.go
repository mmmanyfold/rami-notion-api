package notion

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
)

// request per second to notion api limit
const rate uint64 = 3

var database = map[string]notionapi.DatabaseID{
	"projects":    notionapi.DatabaseID("bee593efdc654282911f3dc5550e144a"),
	"homepage":    notionapi.DatabaseID("a79aece399014bc282a27024de23464a"),
	"transcripts": notionapi.DatabaseID("d815aa37777a4b04812f38b0b9d81b89"),
}

func GetHomePageAssets(client *notionapi.Client) (assets []rami.HomePageAsset, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), database["homepage"], &dbRequest)
	if err != nil {
		return assets, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			if len(r.Properties["File"].(*notionapi.FilesProperty).Files) > 0 {
				asset := rami.HomePageAsset{
					UUID: string(r.ID),
				}
				if len(r.Properties["File"].(*notionapi.FilesProperty).Files) > 0 {
					for _, f := range r.Properties["File"].(*notionapi.FilesProperty).Files {
						asset.Files = append(asset.Files, f)
					}
				}
				assets = append(assets, asset)
			}

		}
	}

	return assets, nil
}

func GetProjects(client *notionapi.Client, assets []rami.HomePageAsset) error {
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

	// take one from the request stack for the first request to retrieve all the ids
	_, _, err = rateLimiter.Take()
	if err != nil {
		return err
	}

	db, err := client.Database.Query(context.Background(), database["projects"], &dbRequest)
	if err != nil {
		return err
	}

	if len(db.Results) > 0 {
		var projects []rami.Project
		for _, r := range db.Results {
			projects = append(projects, rami.Project{
				UUID:           string(r.ID),
				ID:             ProcessRichTextProperty(&r, "ID"),
				Title:          r.Properties["Title"].(*notionapi.TitleProperty).Title[0].Text.Content,
				Tags:           ProcessTags(&r),
				Year:           ProcessYears(&r),
				Thumbnail:      ProcessThumbnail(&r),
				Medium:         ProcessRichTextProperty(&r, "Medium"),
				Description:    ProcessRichTextProperty(&r, "Description"),
				HomePageAssets: processHomePageAsset(&r, assets),
			})
		}
		//fmt.Printf("%+v\n", projects)
		//fmt.Printf("len %d\n", len(projects))
	}

	if err := rateLimiter.Close(); err != nil {
		return err
	}

	return nil
}

func ProcessTags(page *notionapi.Page) (tags []rami.Tag) {
	if multiSelectProperty, ok := page.Properties["Tags"].(*notionapi.MultiSelectProperty); ok {
		for _, t := range multiSelectProperty.MultiSelect {
			tags = append(tags, rami.Tag(t.Name))
		}
	}
	return tags
}

func ProcessYears(page *notionapi.Page) (year string) {
	if selectProperty, ok := page.Properties["Year"].(*notionapi.SelectProperty); ok {
		year = selectProperty.Select.Name
	}
	return year
}

func ProcessThumbnail(page *notionapi.Page) (thumbnailUrl string) {
	if filesProperty, ok := page.Properties["Thumbnail"].(*notionapi.FilesProperty); ok {
		thumbnailUrl = filesProperty.Files[0].File.URL
	}

	return thumbnailUrl
}

func ProcessRichTextProperty(page *notionapi.Page, field string) (text string) {
	if textProperty, ok := page.Properties[field].(*notionapi.RichTextProperty); ok {
		for _, rt := range textProperty.RichText {
			if len(textProperty.RichText) > 0 {
				text += rt.Text.Content
			}
		}
	}

	return text
}

func processHomePageAsset(page *notionapi.Page, assets []rami.HomePageAsset) (homePageAssets rami.HomePageAsset) {
	if relationProperty, ok := page.Properties["Homepage Assets"].(*notionapi.RelationProperty); ok {
		for _, asset := range assets {
			if asset.UUID == string(relationProperty.Relation[0].ID) {
				homePageAssets = rami.HomePageAsset{
					UUID:  asset.UUID,
					Files: asset.Files,
				}
			}
		}
	}

	return homePageAssets
}
