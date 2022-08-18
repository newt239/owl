package functions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

func GetNarou(discord *discordgo.Session) {
	narouChannel, _ := discord.Channel(GetConfig(discord, "NAROU_CHANNEL"))
	message, _ := discord.ChannelMessage(narouChannel.ID, narouChannel.LastMessageID)
	lastUrl := message.Content

	c := colly.NewCollector()
	flag := false
	var newLinkList []string
	c.OnHTML(".subtitle", func(e *colly.HTMLElement) {
		link, _ := e.DOM.Find("a").Attr("href")
		if strings.Split(lastUrl, ".com")[1] == link || flag {
			flag = true
			newLinkList = append(newLinkList, "https://ncode.syosetu.com"+link)
		}
	})
	c.Visit("https://ncode.syosetu.com/n2267be/")

	for i := 0; i < len(newLinkList); i++ {
		discord.ChannelMessageSend(narouChannel.ID, newLinkList[i])
	}
}
