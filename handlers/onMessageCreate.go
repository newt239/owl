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
		if strings.Contains(m.Content, "https://") {
			client := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_KEY")))
			database, _ := client.Database.Get(context.Background(), notionapi.DatabaseID(os.Getenv("MUSIC_DATABASE_ID")))
			arr := strings.Split(m.Content, " ")
			if strings.Contains(arr[0], "youtube.com") || strings.Contains(arr[0], "youtu.be") {
				service, _ := youtube.NewService(context.Background(), option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
				var videoId string
				if strings.Contains(arr[0], "youtube.com") {
					videoId = strings.Split(arr[0], "?v=")[1]
				} else if strings.Contains(arr[0], "youtu.be") {
					videoId = strings.Split(arr[0], "/")[2]
				}
				videoList, _ := service.Videos.List([]string{"snippet"}).Id(videoId).Do()
				video := videoList.Items[0]
				descriptionParagraph := []notionapi.Block{}
				for _, v := range strings.Split(video.Snippet.Description, "\n\n") {
					descriptionParagraph = append(descriptionParagraph, notionapi.ParagraphBlock{
						BasicBlock: notionapi.BasicBlock{
							Object: "block",
							Type:   "paragraph",
						},
						Paragraph: notionapi.Paragraph{
							RichText: []notionapi.RichText{
								{Text: notionapi.Text{Content: v + "\n"}},
							},
						},
					})
				}
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
					Children: descriptionParagraph,
					Icon: &notionapi.Icon{
						Type: "external",
						External: &notionapi.FileObject{
							URL: "https://www.youtube.com/s/desktop/e06db45c/img/favicon_144x144.png",
						},
					},
					Cover: &notionapi.Image{
						Type: "external",
						External: &notionapi.FileObject{
							URL: getLargerYoutubeThambnail(*video.Snippet.Thumbnails),
						},
					},
				}
				page, err := client.Page.Create(context.Background(), &requestBody)
				if err != nil {
					fmt.Println(err)
				}
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title: "Add to Music DB in Notion",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Link",
							Value: page.URL,
						},
					},
					Image: &discordgo.MessageEmbedImage{
						URL: getLargerYoutubeThambnail(*video.Snippet.Thumbnails),
					},
					// https://gist.github.com/thomasbnt/b6f455e2c7d743b796917fa3c205f812
					Color: 16777215,
				})
			}
		}
	}
}

func getLargerYoutubeThambnail(thumbnails youtube.ThumbnailDetails) string {
	if thumbnails.Maxres != nil {
		return thumbnails.Maxres.Url
	} else if thumbnails.Standard != nil {
		return thumbnails.Standard.Url
	} else if thumbnails.High != nil {
		return thumbnails.High.Url
	} else {
		return thumbnails.Default.Url
	}
}
