package main

import (
	"fmt"
	"os"

	bs "github.com/csduarte/mattermost-integration/botserver"
	"github.com/paddycarey/gophy"
)

func main() {
	server, err := bs.NewServer("./config.json")
	if err != nil {
		fmt.Println("Initialization error -", err)
		os.Exit(1)
	}

	// Handles all mattermost incoming
	server.HandleAllChatter("help", sendHelp)

	// Handles all timmer triggers, not yet implemented
	// server.HandleAllTimmers(imAlive)

	// Handles all external triggers, not yet implemented
	// server.HandleAllExternal(imAlive)

	server.HandleChatter("giphy", ".*", getGiphyImage)
	// server.HandleExternal("weatherbot", getWeather)
	// server.HandleTimmer("timmmer", getWeather)

	server.HandleChatter("homer", "tell.*donuts", donutInfo)
	server.HandleChatter("homer", "tell.*family.*", familyInfo)
	server.HandleChatter("homer", ".*", homerMissing)
	server.HandleSomeChatter([]{"home", "giphy"}, ".*", sendHelp)
	server.HandleAllChatter(".*", sendHelp)

	err = server.Start()
	if err != nil {
		panic(err)
	}
}

func sendHelp(context *bs.Context) {
	context.SetMessage("Here is your help!")
	context.SetIconURL("https://media2.giphy.com/media/snJ4LpyvG7OYE/200w.gif#2")
}

func imAlive(context *bs.Context) {
	context.SetMessage("Pong")
}

func getGiphyImage(context *bs.Context) {
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
		r.Send()
	} else {
		context.SetMessage("Sorry, pal.")
	}
}

func getWeather(context *bs.Context) {
	context.SetMessage("It's normal")
}

func getSpecialWeather(context *bs.Context) {
	context.SetMessage("It's not normal")
}

func homerMissing(context *bs.Context) {
	context.SetIconURL("https://media1.giphy.com/media/sTUWqCKtxd01W/200.gif#7")
	context.SetMessage("D'oh! I have no idea what you're talking about. I am so smrt")
}

func donutInfo(context *bs.Context) {
	r := context.SeparateResponse()
	r.SetUsername("Homer Simpson")
	r.SetIconURL("https://media1.giphy.com/media/sTUWqCKtxd01W/200.gif#7")
	r.AttachmentTitle("My Favorite Donut Title")
	r.AttachmentTitleLink("https://en.wikipedia.org/wiki/Doughnut")
	r.AttachmentImageURL("http://superawesomevectors.com/wp-content/uploads/2014/03/free-vector-donut-drawing-800x565.jpg")
	r.AttachmentAuthorIcon("http://vignette3.wikia.nocookie.net/simpsons/images/b/b0/HomerSimpson5.gif/revision/latest?cb=20141025153213")
	r.AttachmentAuthorLink("https://en.wikipedia.org/wiki/Homer_Simpson")
	r.AttachmentAuthorName("Homer")
	r.AttachmentColor("#fe4ea1")
	r.AttachmentFallback("This is my fallback message when the client doesn't support attachments -- about donuts")
	r.AttachmentPretext("The message sent before the wonderful donut attachment")
	r.AttachmentText("This is the actual message about donuts")
	r.AttachmentAddField("Full Row", "A really really long message about donuts", false)
	r.AttachmentAddField("Half Row", "A shorter message about donuts", true)
	r.AttachmentAddField("Other Half Row", "Another short message about donuts", true)
	r.Send()
}

func familyInfo(context *bs.Context) {
	r := context.SeparateResponse()
	r.SetUsername("Homer Simpson")
	r.SetIconURL("https://media1.giphy.com/media/sTUWqCKtxd01W/200.gif#7")
	r.AttachmentTitle("My Stupid Family")
	r.AttachmentImageURL("http://assets2.ignimgs.com/2014/10/01/the-simpsons-couch-1280jpg-552cbc_1280w.jpg")
	r.Send()
}
