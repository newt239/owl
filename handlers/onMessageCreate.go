package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv("CLIENT_ID")
	u := m.Author
	if u.ID != clientId {
		if strings.HasPrefix(m.Content, "https://") {
			client := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_KEY")))
			database, _ := client.Database.Get(context.Background(), notionapi.DatabaseID(os.Getenv("MUSIC_DATABASE_ID")))
			arr := strings.Split(m.Content, " ")

			service, _ := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
			videoList, _ := service.Videos.List([]string{"snippet"}).Id(strings.Split(arr[0], "?v=")[1]).Do()
			video := videoList.Items[0]
			fmt.Println(video.Snippet.Description)
			requestBody := notionapi.PageCreateRequest{
				Parent: notionapi.Parent{
					Type:       notionapi.ParentTypeDatabaseID,
					DatabaseID: notionapi.DatabaseID(database.ID),
				},
				Properties: notionapi.Properties{
					"Name": notionapi.TitleProperty{
						Title: []notionapi.RichText{
							{Text: notionapi.Text{Content: video.Snippet.Title}},
						},
					},
					"URL": notionapi.URLProperty{
						URL: arr[0],
					},
				},
				Children: []notionapi.Block{
					notionapi.ParagraphBlock{
						BasicBlock: notionapi.BasicBlock{
							Object: "block",
							Type:   "paragraph",
						},
						Paragraph: notionapi.Paragraph{
							RichText: []notionapi.RichText{
								{Text: notionapi.Text{Content: video.Snippet.Description}},
							},
						},
					},
				},
				Icon: &notionapi.Icon{
					Type: "external",
					External: &notionapi.FileObject{
						URL: "https://www.youtube.com/s/desktop/e06db45c/img/favicon_144x144.png",
					},
				},
				Cover: &notionapi.Image{
					Type: "external",
					External: &notionapi.FileObject{
						URL: video.Snippet.Thumbnails.Maxres.Url,
					},
				},
			}
			page, err := client.Page.Create(context.Background(), &requestBody)
			if err != nil {
				fmt.Println(err)
			}
			s.ChannelMessageSend(m.ChannelID, page.URL)
		}
	}
}
