package handlers

import (
	"context"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jomei/notionapi"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv("CLIENT_ID")
	u := m.Author
	if u.ID != clientId {
		client := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_KEY")))
		database, _ := client.Database.Get(context.Background(), notionapi.DatabaseID(os.Getenv("MUSIC_DATABASE_ID")))

		_, err := s.ChannelMessageSend(m.ChannelID, database.URL)
		if err != nil {
			log.Println("Error sending message: ", err)
		}
	}
}
