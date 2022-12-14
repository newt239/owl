package functions

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type WeatherResponseStruct []struct {
	PublishingOffice string    `json:"publishingOffice"`
	ReportDatetime   time.Time `json:"reportDatetime"`
	TimeSeries       []struct {
		TimeDefines []time.Time `json:"timeDefines"`
		Areas       []struct {
			Area struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"area"`
			Pops         []string
			WeatherCodes []string `json:"weatherCodes"`
			Weathers     []string `json:"weathers"`
			Winds        []string `json:"winds"`
		} `json:"areas"`
	} `json:"timeSeries"`
	TempAverage struct {
		Areas []struct {
			Area struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"area"`
			Min string `json:"min"`
			Max string `json:"max"`
		} `json:"areas"`
	} `json:"tempAverage,omitempty"`
	PrecipAverage struct {
		Areas []struct {
			Area struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"area"`
			Min string `json:"min"`
			Max string `json:"max"`
		} `json:"areas"`
	} `json:"precipAverage,omitempty"`
}

func GetWeather(discord *discordgo.Session) {
	weatherChannel, _ := discord.Channel(GetConfig(discord, "WEATHER_CHANNEL"))

	url := "https://www.jma.go.jp/bosai/forecast/data/forecast/110000.json"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	raw, _ := client.Do(req)
	body, _ := io.ReadAll(raw.Body)
	var response WeatherResponseStruct
	json.Unmarshal(body, &response)

	pops := response[0].TimeSeries[1].Areas[1].Pops
	timeDefines := response[0].TimeSeries[1].TimeDefines

	if strings.Contains(response[0].TimeSeries[0].Areas[1].Weathers[0], "雨") {
		title := "埼玉県南部の天気 - " + strconv.Itoa(response[0].ReportDatetime.Hour()) + "時発表\n"
		day1Weather := strings.ReplaceAll(response[0].TimeSeries[0].Areas[1].Weathers[0], "晴れ", "🌞晴れ")
		day1Weather = strings.ReplaceAll(day1Weather, "くもり", "☁くもり")
		day1Weather = strings.ReplaceAll(day1Weather, "雨", "☔雨")
		day1Weather = strings.ReplaceAll(day1Weather, "雷", "⚡雷")
		body := "`" + strconv.Itoa(response[0].TimeSeries[0].TimeDefines[0].Day()) + "日:`" + day1Weather + " \n"
		for i := 0; i < len(pops); i++ {
			if i == 0 {
				body += "\n> 降水確率\n"
			}
			weatherCount, _ := strconv.Atoi(pops[i])
			icon := strings.Repeat("🌧", weatherCount/10) + strings.Repeat("➖", 10-weatherCount/10)
			hour := timeDefines[i].Hour()
			var hourStr string
			if hour < 10 {
				hourStr = "0" + strconv.Itoa(hour)
			} else {
				hourStr = strconv.Itoa(hour)
			}
			body += "`" + hourStr + "時` " + icon + " " + pops[i] + "%\n"
		}
		discord.ChannelMessageSendEmbed(weatherChannel.ID, &discordgo.MessageEmbed{
			Title:       title,
			Description: body,
			Color:       1752220,
		})
	}
}
