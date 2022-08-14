package functions

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

func GetShipNews(discord *discordgo.Session) {
	configChannel, _ := discord.Channel(os.Getenv("CHANNEL_ID"))
	configMessage, _ := discord.ChannelMessage(configChannel.ID, configChannel.LastMessageID)
	channelId := strings.Split(configMessage.Content, "\n")[0]
	shnewsChannel, _ := discord.Channel(channelId)
	message, _ := discord.ChannelMessage(channelId, shnewsChannel.LastMessageID)
	lastUrl := message.Content

	c := colly.NewCollector()
	flag := true
	var newLinkList []string
	c.OnHTML(".index-list", func(e *colly.HTMLElement) {
		link, _ := e.DOM.Find("a").Attr("href")
		link = "https://www.sakaehigashi.ed.jp" + link
		if link != lastUrl && flag {
			newLinkList = append(newLinkList, link)
		} else {
			flag = false
		}
	})
	c.Visit("http://www.sakaehigashi.ed.jp/news/")

	for i := len(newLinkList) - 1; i >= 0; i-- {
		discord.ChannelMessageSend(shnewsChannel.ID, newLinkList[i])
	}
}
