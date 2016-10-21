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
	bot.HandleAll("ping", imAlive)

	bot.HandleOne("giphy", ".*", getGiphyImage)
	bot.HandleOne("weatherbot", ".*", getWeather)
	bot.HandleOne("mint", "basic", testBasic)
	bot.HandleOne("mint", "full", testFull)

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
	fmt.Println("ignore:", gifs[0].URL)

	if total > 0 {
		r := context.SeparateResponse()
		r.SetMessage("Boohyah Grandma")
		// r.AttachmentImageURL("http://giphy.com/gifs/cute-dog-running-13mLwGra9bNEKQ")
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

func testBasic(context *s.Context) {
	context.SetMessage("Hi you")
}

func testFull(context *s.Context) {
	context.SetMessage("Hi me")
}
