package functions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

func GetShipNews(discord *discordgo.Session) {
	shnewsChannel, _ := discord.Channel(GetConfig(discord, "SHNEWS_CHANNEL"))
	message, _ := discord.ChannelMessage(shnewsChannel.ID, shnewsChannel.LastMessageID)
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
