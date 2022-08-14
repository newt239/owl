package functions

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

func GetNarou(discord *discordgo.Session) {
	configChannel, _ := discord.Channel(os.Getenv("CHANNEL_ID"))
	configMessage, _ := discord.ChannelMessage(configChannel.ID, configChannel.LastMessageID)
	channelId := strings.Split(configMessage.Content, "\n")[1]
	shnewsChannel, _ := discord.Channel(channelId)
	message, _ := discord.ChannelMessage(channelId, shnewsChannel.LastMessageID)
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
		discord.ChannelMessageSend(shnewsChannel.ID, newLinkList[i])
	}
}
