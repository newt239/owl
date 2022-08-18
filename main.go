package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/newt239/owl/functions"
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

	defer discord.Close()

	fmt.Println("owl is running")

	functions.GetShipNews(discord)
	functions.GetNarou(discord)

	ticker := time.NewTicker(time.Hour)
	fmt.Println("タイマーを開始")
	go func() {
		for t := range ticker.C {
			fmt.Println(t)
			functions.GetShipNews(discord)
			functions.GetNarou(discord)
			if t.Hour() == 5 {
				functions.GetWeather(discord)
			}
		}
	}()

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stopBot
}
