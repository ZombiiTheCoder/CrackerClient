package main

import (
	"embed"
	"math/rand"
	"os"
	"strconv"
	"time"

	webview "github.com/webview/webview_go"
)

var (
	edition string

	LocalAppData = os.Getenv("localappdata")
	RoamingAppData = os.Getenv("appdata")
	//go:embed www
	www embed.FS
	//go:embed data
	data embed.FS

	random = rand.New(rand.NewSource(time.Hour.Milliseconds()))
	port = random.Intn(9999)
)

func main() {
	
	w := webview.New(edition == "dev")
	defer w.Destroy()
	if (edition != "dev") { edition = ""; w.SetTitle("Cracker Client") } else { w.SetTitle("Cracker Client - Dev") }
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:"+strconv.Itoa(port))
	setResourceIcon(w, "CLIENTLOGO")
	bindFunctions(w)
	go serverEmbed(www, "localhost", strconv.Itoa(port))
	w.Run()
}