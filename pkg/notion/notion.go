package notion

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

var Databases = map[string]notionapi.DatabaseID{
	"projects":                      notionapi.DatabaseID("bee593efdc654282911f3dc5550e144a"), // resource
	"homepage":                      notionapi.DatabaseID("a79aece399014bc282a27024de23464a"),
	"transcripts":                   notionapi.DatabaseID("d815aa37777a4b04812f38b0b9d81b89"),
	"info":                          notionapi.DatabaseID("74db7bbeb10b41dca217d55e9a675e3e"), // resource
	"cv-exhibitions-and-screenings": notionapi.DatabaseID("9090f44a583049d7a7fa478b2dd329a8"), // resource
	"cv-additional":                 notionapi.DatabaseID("78403af9f31145ce98c7a9ffa57931f8"), // resource
}

func GetInfoDB(client *notionapi.Client) (rows []rami.Info, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), Databases["info"], &dbRequest)
	if err != nil {
		return rows, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			var row rami.Info
			row.UUID = string(r.ID)
			row.Tag = processSelectProperty(&r, "Tag")
			row.Line1 = processRichTextProperty(&r, "Line 1")
			row.Line2 = processRichTextProperty(&r, "Line 2")
			row.Line3 = processRichTextProperty(&r, "Line 3")
			row.Line4 = processRichTextProperty(&r, "Line 4")
			row.URL = processTitleProperty(&r, "URL")
			row.Download = processFilesProperty(&r, "Download")
			rows = append(rows, row)
		}
	}

	return rows, nil
}

func GetCVExhibitionsAndScreeningDB(client *notionapi.Client) (rows []rami.CVExhibitionsAndScreening, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), Databases["cv-exhibitions-and-screenings"], &dbRequest)
	if err != nil {
		return rows, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			var row rami.CVExhibitionsAndScreening
			row.UUID = string(r.ID)
			row.Title = processRichTextProperty(&r, "Title")
			row.Description = processRichTextProperty(&r, "Description")
			row.Extra = processRichTextProperty(&r, "Extra")
			row.URL = processTitleProperty(&r, "URL")
			row.Download = processFilesProperty(&r, "Download")
			row.Year = processSelect(&r, "Year")
			// TODO: For Project Press page relation prop
			rows = append(rows, row)
		}
	}

	return rows, err
}

func GetCVAdditionalDB(client *notionapi.Client) (rows []rami.CVAdditional, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), Databases["cv-additional"], &dbRequest)
	if err != nil {
		return rows, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			var row rami.CVAdditional
			row.UUID = string(r.ID)
			row.Title = processRichTextProperty(&r, "Title")
			row.Description = processRichTextProperty(&r, "Description")
			row.Detail = processRichTextProperty(&r, "Detail")
			row.URL = processTitleProperty(&r, "URL")
			row.Tag = processSelectProperty(&r, "Tag")
			row.Download = processFilesProperty(&r, "Download")
			// TODO: For Project Press page relation prop
			rows = append(rows, row)
		}
	}

	return rows, err
}

func GetTranscripts(client *notionapi.Client) (transcripts []rami.Transcript, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), Databases["transcripts"], &dbRequest)
	if err != nil {
		return transcripts, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			var transcript rami.Transcript
			transcript.UUID = string(r.ID)
			if projectRelationProperty, ok := r.Properties["Project"].(*notionapi.RelationProperty); ok {
				transcript.ProjectUUID = string(projectRelationProperty.ID)
				pagination := notionapi.Pagination{
					StartCursor: "",
					PageSize:    0,
				}
				pageBlocks, err := client.Block.GetChildren(context.TODO(), notionapi.BlockID(transcript.UUID), &pagination)
				if err != nil {
					errMessage := fmt.Sprintf("failed to get transcript blocks for page id: %s", transcript.UUID)
					fmt.Println(errors.Wrap(err, errMessage))
					break
				}
				transcript.Blocks = pageBlocks.Results
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

	db, err := client.Database.Query(context.Background(), Databases["homepage"], &dbRequest)
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
						asset.NotionFiles = append(asset.NotionFiles, f)
					}
				}
				assets = append(assets, asset)
			}
		}
	}

	return assets, nil
}

func GetProjects(client *notionapi.Client, assets []rami.HomePageAsset, transcripts []rami.Transcript) (rows []rami.Project, err error) {
	dbRequest := notionapi.DatabaseQueryRequest{
		Filter:      nil,
		Sorts:       nil,
		StartCursor: "",
		PageSize:    0,
	}

	db, err := client.Database.Query(context.Background(), Databases["projects"], &dbRequest)
	if err != nil {
		return rows, err
	}

	if len(db.Results) > 0 {
		for _, r := range db.Results {
			title := processTitleProperty(&r, "Title")
			id := processRichTextProperty(&r, "ID")
			rows = append(rows, rami.Project{
				UUID:           string(r.ID),
				ID:             id,
				Title:          title,
				Tags:           processTags(&r),
				Year:           processSelect(&r, "Year"),
				Thumbnail:      processThumbnail(&r),
				Medium:         processRichTextProperty(&r, "Medium"),
				Description:    processRichTextProperty(&r, "Description"),
				HomePageAssets: processHomePageAsset(&r, assets),
				Transcript:     processTranscript(&r, transcripts),
				Slug:           id + "-" + Slug(title),
			})
		}
	}

	return rows, nil
}

func processTitleProperty(page *notionapi.Page, fieldName string) (title string) {
	if titleProperty, ok := page.Properties[fieldName].(*notionapi.TitleProperty); ok {
		if len(titleProperty.Title) > 0 {
			title = titleProperty.Title[0].Text.Content
		}
	}

	return title
}

func processTags(page *notionapi.Page) (tags []rami.Tag) {
	if multiSelectProperty, ok := page.Properties["Tags"].(*notionapi.MultiSelectProperty); ok {
		for _, t := range multiSelectProperty.MultiSelect {
			lowerT := strings.ToLower(t.Name)
			tags = append(tags, rami.Tag(lowerT))
		}
	}
	return tags
}

func processSelectProperty(page *notionapi.Page, fieldName string) (tag string) {
	if selectProperty, ok := page.Properties[fieldName].(*notionapi.SelectProperty); ok {
		tag = selectProperty.Select.Name
	}
	return tag
}

func processSelect(page *notionapi.Page, fieldName string) (year string) {
	if selectProperty, ok := page.Properties[fieldName].(*notionapi.SelectProperty); ok {
		year = selectProperty.Select.Name
	}
	return year
}

func processThumbnail(page *notionapi.Page) (thumbnailUrl string) {
	if filesProperty, ok := page.Properties["Thumbnail"].(*notionapi.FilesProperty); ok {
		thumbnailUrl, _ = url.PathUnescape(filesProperty.Files[0].File.URL)
	}

	return thumbnailUrl
}

func processRichTextProperty(page *notionapi.Page, fieldName string) (text string) {
	if textProperty, ok := page.Properties[fieldName].(*notionapi.RichTextProperty); ok {
		for _, rt := range textProperty.RichText {
			if len(textProperty.RichText) > 0 {
				text += rt.Text.Content
			}
		}
	}

	return text
}

func processFilesProperty(page *notionapi.Page, fieldName string) (files []rami.File) {
	if filesProperty, ok := page.Properties[fieldName].(*notionapi.FilesProperty); ok {
		for _, rt := range filesProperty.Files {
			if len(filesProperty.Files) > 0 {
				file := rami.File{
					Url: rt.File.URL,
				}
				files = append(files, file)
			}
		}
	}

	return files
}

func processHomePageAsset(page *notionapi.Page, assets []rami.HomePageAsset) (homePageAssets rami.HomePageAsset) {
	if relationProperty, ok := page.Properties["Homepage Assets"].(*notionapi.RelationProperty); ok {
		for _, asset := range assets {
			if asset.UUID == string(relationProperty.Relation[0].ID) {
				var files []rami.File
				for _, u := range asset.NotionFiles {
					s3Url := u.File.URL
					files = append(files, rami.File{
						Url: s3Url,
					})
				}
				homePageAssets = rami.HomePageAsset{
					UUID:  asset.UUID,
					Files: files,
					Type:  asset.Type,
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
