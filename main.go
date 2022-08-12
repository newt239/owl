package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/newt239/owl/handlers"
)

func main() {
	godotenv.Load(".env")

	discord, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("ログインに失敗しました")
		fmt.Println(err)
	}

	discord.AddHandler(handlers.OnMessageCreate)
	discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	// 直近の関数（main）の最後に実行される
	defer discord.Close()

	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot
}
