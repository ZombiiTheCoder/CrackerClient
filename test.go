package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	webview "github.com/webview/webview_go"
)

func main() {
	
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Minecraft Login")
	w.SetSize(800, 600, webview.HintNone)
	// random := rand.New(rand.NewSource(time.Hour.Milliseconds()))
	// port := random.Intn(9999);
	w.Navigate("https://sisu.xboxlive.com/connect/XboxLive/?state=login&cobrandId=8058f65d-ce06-4c30-9559-473c9275a65d&tid=896928775&ru=https%3A%2F%2Fwww.minecraft.net%2Fen-us%2Flogin&aid=1142970254")
	w.Bind("println", func (fn any)  {
		fmt.Println(fn)
	})
	w.Bind("save", func (data string)  {
		os.WriteFile("cookie.txt", []byte(data), fs.FileMode(os.O_CREATE))
	})
	w.Bind("close_window", func() {
		w.Terminate()
	})
	go func() {
		for {
			w.Dispatch(func() {
				w.Eval("setTimeout(()=>{if (document.domain == 'www.minecraft.net') { save(document.cookie); close_window() }}, 8000)")
			})
			time.Sleep(1 * time.Second)
		}
	}()
	// w.Navigate("http://localhost:"+strconv.Itoa(port))
	// setResourceIcon(w, "CLIENTLOGO")
	// bindFunctions(w)
	// go serverEmbed(www, "localhost", strconv.Itoa(port))
	w.Run()
}