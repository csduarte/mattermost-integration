package main

import (
	"fmt"
	"os"

	s "github.com/csduarte/mattermost-integration/server"
)

const (
	server1 = "mattermost-staging-generic"
	server2 = "mattermost-production"
)

func main() {
	bot, err := s.NewIntegrationServer("./config.json")
	if err != nil {
		fmt.Println("Startup error -", err)
		os.Exit(1)
	}

	bot.HandleAll("me", getGiphyImage)
	bot.HandleSome([]string{server1, server2}, "weather", getWeather)
	bot.HandleOne(server1, "traffic", getSpecialWeather)

	err = bot.Start()
	if err != nil {
		panic(err)
	}
}

func getGiphyImage(context s.Context) {
	// get giphy image
	context.AddAttachment("imageData")
	context.Respond("Image")
}

func getWeather(context s.Context) {
	// do some work
	context.Respond("The weather is nice, but getting colder")
}

func getSpecialWeather(context s.Context) {
	// do some work
	context.Respond("The weather is especially nice, but getting warmer")
}
