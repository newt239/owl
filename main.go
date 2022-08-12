package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("ログインに失敗しました")
		fmt.Println(err)
	}

	discord.AddHandler(onMessageCreate)
	discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	// 直近の関数（main）の最後に実行される
	defer discord.Close()

	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	clientId := os.Getenv("CLIENT_ID")
	u := m.Author
	fmt.Printf("%20s %20s(%20s) > %s\n", m.ChannelID, u.Username, u.ID, m.Content)
	if u.ID != clientId {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hi!")
		if err != nil {
			log.Println("Error sending message: ", err)
		}
	}
}
