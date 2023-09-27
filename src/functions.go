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
    bind := func(name string, fn interface{}) {
        if err := w.Bind(name, fn); err != nil {
            fmt.Printf("Failed to bind %s: %v\n", name, err)
        }
    }

    // Auth
    bind("login_with_microsoft", func ()  { login_with_microsoft(w) })

    bind("print", func(str any) { fmt.Println(str) })

    fileExist := func(filePath string) bool {
        _, err := os.Stat(filePath)
        return err == nil
    }
    embedFileExist := func(filePath string) bool {
        fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
        _, fileName := path.Split(filePath)
        _, err := fs.Stat(fileSys, fileName)
        return err == nil
    }

    bind("FileExist", fileExist)
    bind("Embed_FileExist", embedFileExist)

    readFile := func(filePath string) string {
        content, _ := os.ReadFile(filePath)
        return string(content)
    }
    embedReadFile := func(filePath string) string {
        fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
        _, fileName := path.Split(filePath)
        content, _ := fs.ReadFile(fileSys, fileName)
        return string(content)
    }

    bind("ReadFile", readFile)
    bind("Embed_ReadFile", embedReadFile)

    readFileAsDataUrl := func(filePath string) string {
        content, _ := os.ReadFile(filePath)
        return dataurl.EncodeBytes(content)
    }
    embedReadFileAsDataUrl := func(filePath string) string {
        fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
        _, fileName := path.Split(filePath)
        content, _ := fs.ReadFile(fileSys, fileName)
        return dataurl.EncodeBytes(content)
    }

    bind("ReadFileAsDataUrl", readFileAsDataUrl)
    bind("Embed_ReadFileAsDataUrl", embedReadFileAsDataUrl)

    copyFile := func(filePath, newFilePath string) {
        os.MkdirAll(path.Dir(newFilePath), fs.FileMode(os.O_CREATE))
        z, _ := os.ReadFile(filePath)
        os.WriteFile(newFilePath, z, fs.FileMode(os.O_CREATE))
    }
    embedCopyFile := func(filePath, newFilePath string) {
        os.MkdirAll(path.Dir(newFilePath), fs.FileMode(os.O_CREATE))
        fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
        _, fileName := path.Split(filePath)
        content, _ := fs.ReadFile(fileSys, fileName)
        os.WriteFile(newFilePath, content, fs.FileMode(os.O_CREATE))
    }

    bind("CopyFile", copyFile)
    bind("Embed_CopyFile", embedCopyFile)

    bind("RemoveFile", func(filePath string) { os.Remove(filePath) })
    bind("RemoveDir", func(path string) { os.RemoveAll(path) })

    writeFile := func(filePath, content string) {
        os.MkdirAll(path.Dir(filePath), fs.FileMode(os.O_CREATE))
        q, _ := os.Create(filePath)
        q.WriteString(content)
        q.Close()
    }

    bind("WriteFile", writeFile)
    bind("LocalAppdata", func() string { return LocalAppData })
    bind("RoamingAppdata", func() string { return RoamingAppData })
    bind("edition", func() string { return edition })
    bind("port", func() int { return port; })
    bind("execute", func(cwd string, prg string, args ...string) {
        cmd := exec.Command(prg, args...)
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
        // Realtime std::out logging
        go func() {
            scanner := bufio.NewScanner(stdout)
            for scanner.Scan() {
                m := scanner.Text()
                fmt.Printf("%v\n", m)
                w.Dispatch(func() {
                    // Log to HTML logger
                    w.Eval(fmt.Sprintf(`
z3 = document.createElement("pre"); z3.innerText = "%v"; z3.setAttribute("class", "log");
document.getElementsByClassName("console")[0].append(z3)
if (Math.round(document.getElementsByClassName("console")[0].scrollTop) >= (document.getElementsByClassName("console")[0].scrollHeight - 859)) {
    document.getElementsByClassName("console")[0].scroll(0, document.getElementsByClassName("console")[0].scrollHeight)
}
`, m))
                })
            }
            cmd.Wait()
        }()

    })
}
