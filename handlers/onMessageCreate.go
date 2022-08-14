package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
	"github.com/newt239/owl/functions"
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
							URL: "https://youtube.com/watch?v=" + videoId,
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
			} else if strings.Contains(m.Content, "https://discord.com/channels/") {
				channelId := strings.Split(m.Content, "/")[5]
				messageId := strings.Split(m.Content, "/")[6]
				message, _ := s.ChannelMessage(channelId, messageId)
				channel, _ := s.Channel(channelId)
				result, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Description: message.Content,
					Color:       15158332,
					Timestamp:   message.Timestamp.Format(time.RFC3339),
					Author: &discordgo.MessageEmbedAuthor{
						Name:    message.Author.Username,
						IconURL: message.Author.AvatarURL(""),
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: channel.Name,
					},
				})
				fmt.Println(result, err)
			}
		}
		if m.Content == "weather" || m.Content == "天気" {
			functions.GetWeather(s)
		}
		if m.Content == "narou" {
			functions.GetNarou(s)
		}
		if m.Content == "ship" {
			functions.GetShipNews(s)
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
