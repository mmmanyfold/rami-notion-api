package notion

import (
	"context"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
	"net/url"
	"time"
)

// Rate request per second to notion API
const Rate uint64 = 3
const Limit = time.Second / 3

var database = map[string]notionapi.DatabaseID{
	"projects":    notionapi.DatabaseID("bee593efdc654282911f3dc5550e144a"),
	"homepage":    notionapi.DatabaseID("a79aece399014bc282a27024de23464a"),
	"transcripts": notionapi.DatabaseID("d815aa37777a4b04812f38b0b9d81b89"),
}

func GetTranscripts(client *notionapi.Client) (transcripts []rami.Transcript, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), database["transcripts"], &dbRequest)
	if err != nil {
		return transcripts, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			var transcript rami.Transcript
			transcript.UUID = string(r.ID)
			if projectRelationProperty, ok := r.Properties["Project"].(*notionapi.RelationProperty); ok {
				transcript.ProjectUUID = string(projectRelationProperty.ID)
			}
			transcripts = append(transcripts, transcript)
		}
	}

	return transcripts, nil
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
			asset := rami.HomePageAsset{
				UUID: string(r.ID),
			}
			if selectProperty, ok := r.Properties["File Type"].(*notionapi.SelectProperty); ok {
				asset.Type = selectProperty.Select.Name
			}
			if len(r.Properties["File"].(*notionapi.FilesProperty).Files) > 0 {
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

func GetProjects(client *notionapi.Client, assets []rami.HomePageAsset, transcripts []rami.Transcript) (projectsResponse rami.ProjectsResponse, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), database["projects"], &dbRequest)
	if err != nil {
		return projectsResponse, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			projectsResponse.AllProjects = append(projectsResponse.AllProjects, rami.Project{
				UUID:           string(r.ID),
				ID:             ProcessRichTextProperty(&r, "ID"),
				Title:          r.Properties["Title"].(*notionapi.TitleProperty).Title[0].Text.Content,
				Tags:           ProcessTags(&r),
				Year:           ProcessYears(&r),
				Thumbnail:      ProcessThumbnail(&r),
				Medium:         ProcessRichTextProperty(&r, "Medium"),
				Description:    ProcessRichTextProperty(&r, "Description"),
				HomePageAssets: processHomePageAsset(&r, assets),
				Transcript:     processTranscript(&r, transcripts),
			})
		}
		//fmt.Printf("%+v\n", projects)
		//fmt.Printf("len %d\n", len(projects))
	}

	projectsResponse.LastRefreshed = time.Now().String()
	return projectsResponse, nil
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
		thumbnailUrl, _ = url.PathUnescape(filesProperty.Files[0].File.URL)
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
				var urls []string
				for _, u := range asset.Files {
					urls = append(urls, u.File.URL)
				}
				homePageAssets = rami.HomePageAsset{
					UUID: asset.UUID,
					Urls: urls,
					Type: asset.Type,
				}
			}

		}
	}

	return homePageAssets
}

func processTranscript(page *notionapi.Page, transcripts []rami.Transcript) (transcript rami.Transcript) {
	if relationProperty, ok := page.Properties["Transcript"].(*notionapi.RelationProperty); ok {
		if len(relationProperty.Relation) > 0 {
			for _, t := range transcripts {
				if t.UUID == string(relationProperty.Relation[0].ID) {
					transcript = t
				}
			}
		}
	}

	return transcript
}
