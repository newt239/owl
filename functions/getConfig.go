package functions

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetConfig(discord *discordgo.Session, key string) string {
	configChannel, _ := discord.Channel(os.Getenv("CHANNEL_ID"))
	configMessage, _ := discord.ChannelMessage(configChannel.ID, configChannel.LastMessageID)
	configList := strings.Split(configMessage.Content, "\n")
	result := os.Getenv("CHANNEL_ID")
	for _, v := range configList {
		configPare := strings.Split(v, "=")
		if configPare[0] == key {
			result = configPare[1]
		}
	}
	return result
}
