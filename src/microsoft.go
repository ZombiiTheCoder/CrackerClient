package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"
	"unsafe"

	webview "github.com/webview/webview_go"
)

func login_with_microsoft(b unsafe.Pointer) {
	
	wz := webview.NewWindow(true, b)
	// defer wz.Destroy()
	wz.SetTitle("Minecraft Login")
	wz.SetSize(800, 600, webview.HintNone)
	// random := rand.New(rand.NewSource(time.Hour.Milliseconds()))
	// port := random.Intn(9999);
	wz.Navigate("https://sisu.xboxlive.com/connect/XboxLive/?state=login&cobrandId=8058f65d-ce06-4c30-9559-473c9275a65d&tid=896928775&ru=https%3A%2F%2Fwww.minecraft.net%2Fen-us%2Flogin&aid=1142970254")
	wz.Bind("println", func (fn any)  {
		fmt.Println(fn)
	})
	wz.Bind("save", func (data string)  {
		os.WriteFile("cookie.txt", []byte(data), fs.FileMode(os.O_CREATE))
	})
	wz.Bind("close_window", func() {
		// wz.Terminate()
	})
	go func() {
		for {
			wz.Dispatch(func() {
				wz.Eval("setTimeout(()=>{if (document.domain == 'www.minecraft.net') { save(document.cookie); close_window() }}, 8000)")
			})
			time.Sleep(1 * time.Second)
		}
	}()
	// w.Navigate("http://localhost:"+strconv.Itoa(port))
	// setResourceIcon(w, "CLIENTLOGO")
	// bindFunctions(w)
	// go serverEmbed(www, "localhost", strconv.Itoa(port))
	wz.Run()
}