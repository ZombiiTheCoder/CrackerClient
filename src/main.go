package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/vincent-petithory/dataurl"
	webview "github.com/webview/webview_go"
)

//go:embed www
var www embed.FS

func main() {
	isEmbed := false

	debug := true


	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("CrackerClient")
	w.SetSize(800, 600, webview.HintNone)
	var (
		LocalAppData = os.Getenv("localappdata")
		RoamingAppData = os.Getenv("appdata")
	)
	fmt.Println()
	w.Navigate("http://localhost:2044")
	w.Bind("wwwDev", func() string { return "D:/Projects/CrackerClient/src/www"; })

	w.Bind("print", func(str string) { fmt.Println(str); })

	w.Bind("ReadEmbedFile", func(name string) string { z, _ := www.ReadFile(name); return string(z) })
	w.Bind("ReadEmbedFileAsDataUrl", func(name string) string { z, _ := www.ReadFile(name); return dataurl.EncodeBytes(z); })
	w.Bind("ReadFile", func(name string) string { z, _ := os.ReadFile(name); return string(z); })
	w.Bind("ReadFileAsDataUrl", func(name string) string { z, _ := os.ReadFile(name); return dataurl.EncodeBytes(z); })
	w.Bind("RemoveFile", func(name string) { os.Remove(name); })
	w.Bind("RemoveDir", func(name string) { os.RemoveAll(name); })
	w.Bind("CopyFile", func(name, new string) { z, _ := os.ReadFile(name); os.WriteFile(new, z, fs.FileMode(os.O_CREATE)) })
	w.Bind("WriteFile", func(name, content string) { os.WriteFile(name, []byte(content), fs.FileMode(os.O_CREATE)); })
	w.Bind("LocalAppdata", func() string { return LocalAppData; })
	w.Bind("RoamingAppdata", func() string { return RoamingAppData; })
	w.Bind("isEmbed", func() bool { return isEmbed; })
	if (isEmbed) {
		go serverEmbed(www, "http://localhost", "2044")
	} else {
		go server("D:/Projects/CrackerClient/src/www", "http://localhost", "2044")
	}
	w.Run()
}

// 