package main

import (
	"fmt"
	"os"

	s "github.com/csduarte/mattermost-integration/server"
	"github.com/paddycarey/gophy"
)

func main() {
	bot, err := s.NewIntegrationServer("./config.json")
	if err != nil {
		fmt.Println("Initialization error -", err)
		os.Exit(1)
	}

	bot.HandleAll("help", sendHelp)
	bot.HandleSome([]string{"giphy", "weatherbot"}, "ping", imAlive)

	bot.HandleOne("giphy", ".*", getGiphyImage)
	bot.HandleOne("weatherbot", "inside", getWeather)
	bot.HandleOne("weatherbot", ".*", getSpecialWeather)

	err = bot.Start()
	if err != nil {
		panic(err)
	}
}

func sendHelp(context *s.Context) {
	context.SetMessage("Here is your help!")
	context.SetIconURL("https://media2.giphy.com/media/snJ4LpyvG7OYE/200w.gif#2")
}

func imAlive(context *s.Context) {
	context.SetMessage("Pong")
}

func getGiphyImage(context *s.Context) {
	co := &gophy.ClientOptions{}
	client := gophy.NewClient(co)

	gifs, total, err := client.SearchGifs("dog", "pg-13", 1, 0)
	if err != nil {
		panic(err)
	}

	if total > 0 {
		r := context.SeparateResponse()
		r.SetMessage("Boohyah Grandma")
		r.AddImageURL(gifs[0].URL)
		r.Send()
	} else {
		context.SetMessage("Sorry, pal.")
	}
}

func getWeather(context *s.Context) {
	context.SetMessage("It's normal")
}

func getSpecialWeather(context *s.Context) {
	context.SetMessage("It's not normal")
}
