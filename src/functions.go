package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/vincent-petithory/dataurl"
	webview "github.com/webview/webview_go"
)

func bindFunctions(w webview.WebView) {

	w.Bind("print", func(str any) { fmt.Println(str); })

	w.Bind("FileExist", func(filePath string) bool { 
		_, err := os.Stat(filePath);
		return err == nil
	})
	w.Bind("Embed_FileExist", func(filePath string) bool {
		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)));
		_, fileName := path.Split(filePath);
		_, err := fs.Stat(fileSys, fileName);
		return err == nil
	})
	w.Bind("Embed_ReadFile", func(filePath string) string {
		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)));
		_, fileName := path.Split(filePath);
		content, _ := fs.ReadFile(fileSys, fileName);
		return string(content)
	})
	w.Bind("Embed_ReadFileAsDataUrl", func(filePath string) string {
		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)));
		_, fileName := path.Split(filePath);
		content, _ := fs.ReadFile(fileSys, fileName);
		return dataurl.EncodeBytes(content); 
	})
	w.Bind("Embed_CopyFile", func(filePath, NewFilePath string) {
		os.MkdirAll(path.Dir(NewFilePath), fs.FileMode(os.O_CREATE));
		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)));
		_, fileName := path.Split(filePath);
		content, _ := fs.ReadFile(fileSys, fileName);
		os.WriteFile(NewFilePath, content, fs.FileMode(os.O_CREATE))
	})
	w.Bind("ReadFile", func(filePath string) string {
		content, _ := os.ReadFile(filePath);
		return string(content);
	})
	w.Bind("ReadFileAsDataUrl", func(filePath string) string {
		content, _ := os.ReadFile(filePath);
		return dataurl.EncodeBytes(content);
	})
	w.Bind("RemoveFile", func(filePath string) { os.Remove(filePath); })
	w.Bind("RemoveDir", func(path string) { os.RemoveAll(path); })
	w.Bind("CopyFile", func(filePath, NewFilePath string) {
		os.MkdirAll(path.Dir(NewFilePath), fs.FileMode(os.O_CREATE));
		z, _ := os.ReadFile(filePath);
		os.WriteFile(NewFilePath, z, fs.FileMode(os.O_CREATE))
	})
	w.Bind("WriteFile", func(filePath, content string) {
		os.MkdirAll(path.Dir(filePath), fs.FileMode(os.O_CREATE));
		q, _ := os.Create(filePath);
		q.WriteString(content); q.Close()
	})
	w.Bind("LocalAppdata", func() string { return LocalAppData; })
	w.Bind("RoamingAppdata", func() string { return RoamingAppData; })
	w.Bind("edition", func() string { return edition; })

	w.Bind("execute", func(cwd string, prg string, args ...string) {
		cmd := exec.Command(prg, args...);
		cmd.Dir = cwd
		stdout, err := cmd.StdoutPipe()
		if err != nil { fmt.Println(err) }
		err = cmd.Start()
		fmt.Println("The command is running")
		if err != nil { fmt.Println(err) }
		// realtime std::out logging
		go func(){
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Printf("%v\n", m)
				w.Dispatch(func() {
					// Log to HTML logger
					w.Eval(fmt.Sprintf(`
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
		
	})
}