package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/vincent-petithory/dataurl"
	webview "github.com/webview/webview_go"
)

var edition string

//go:embed www
var www embed.FS

//go:embed data
var data embed.FS

func main() {


	isEmbed := true

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	if (edition != "dev") { edition = ""; w.SetTitle("Cracker Client") } else {
		w.SetTitle("Cracker Client - Dev")
	}
	w.SetSize(800, 600, webview.HintNone)
	var (
		LocalAppData = os.Getenv("localappdata")
		RoamingAppData = os.Getenv("appdata")
	)
	fmt.Println()
	w.Navigate("http://localhost:2044")

	w.Bind("print", func(str any) { fmt.Println(str); })

	w.Bind("FileExist", func(name string) bool { _, err := os.Stat(name); return err == nil})
	w.Bind("Embed_FileExist", func(name string) bool { a, _ := fs.Sub(data, path.Join("data", path.Dir(name))); _, f := path.Split(name); _, err := fs.Stat(a, f); return err == nil})
	w.Bind("Embed_ReadFile", func(name string) string { a, _ := fs.Sub(data, path.Join("data", path.Dir(name))); _, f := path.Split(name); z, _ := fs.ReadFile(a, f); return string(z) })
	w.Bind("Embed_ReadFileAsDataUrl", func(name string) string { a, _ := fs.Sub(data, path.Join("data", path.Dir(name))); _, f := path.Split(name); z, _ := fs.ReadFile(a, f); return dataurl.EncodeBytes(z); })
	w.Bind("Embed_CopyFile", func(name, new string) { os.MkdirAll(path.Dir(new), fs.FileMode(os.O_CREATE)); a, _ := fs.Sub(data, path.Join("data", path.Dir(name))); _, f := path.Split(name); z, _ := fs.ReadFile(a, f); os.WriteFile(new, z, fs.FileMode(os.O_CREATE)) })
	w.Bind("ReadFile", func(name string) string { z, _ := os.ReadFile(name); return string(z); })
	w.Bind("ReadFileAsDataUrl", func(name string) string { z, _ := os.ReadFile(name); return dataurl.EncodeBytes(z); })
	w.Bind("RemoveFile", func(name string) { os.Remove(name); })
	w.Bind("RemoveDir", func(name string) { os.RemoveAll(name); })
	w.Bind("CopyFile", func(name, new string) { os.MkdirAll(path.Dir(new), fs.FileMode(os.O_CREATE)); z, _ := os.ReadFile(name); os.WriteFile(new, z, fs.FileMode(os.O_CREATE)) })
	w.Bind("WriteFile", func(name, content string) { os.MkdirAll(path.Dir(name), fs.FileMode(os.O_CREATE));  fmt.Println(name); q, _ := os.Create(name); q.WriteString(content); q.Close() })
	w.Bind("LocalAppdata", func() string { return LocalAppData; })
	w.Bind("RoamingAppdata", func() string { return RoamingAppData; })
	w.Bind("isEmbed", func() bool { return isEmbed; })
	w.Bind("edition", func() string { return edition; })

	w.Bind("execute", func(cwd string, prg string, args ...string) {
		z := func(wb webview.WebView, cwd, prg string, args ...string){
			cmd := exec.Command(prg, args...);
			cmd.Dir = cwd
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println(err)
			}

			err = cmd.Start()
			fmt.Println("The command is running")
			if err != nil {
				fmt.Println(err)
			}
			
			// print the output of the subprocess
			go func(){
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Printf("%v\n", m)
					wb.Dispatch(func() {
						wb.Eval(fmt.Sprintf(`
z3 = document.createElement("pre"); z3.innerText = "%v"; z3.setAttribute("class", "log");
document.getElementsByClassName("console")[0].append(z3)
if (Math.round(document.getElementsByClassName("console")[0].scrollTop) >= (document.getElementsByClassName("console")[0].scrollHeight - 659)) {
	document.getElementsByClassName("console")[0].scroll(0, document.getElementsByClassName("console")[0].scrollHeight)
}
`, m))
					})
				}
				cmd.Wait()
			}()
			
		}
		z(w, cwd, prg, args...)
	})

	if (isEmbed) {
		go serverEmbed(www, "localhost", "2044")
	} else {
		go server("D:/Projects/CrackerClient/src/www", "localhost", "2044")
	}
	w.Run()
}