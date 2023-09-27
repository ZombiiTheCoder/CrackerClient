package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/twinj/uuid"

	"github.com/vincent-petithory/dataurl"
	webview "github.com/webview/webview_go"
)

func bindFunctions(w webview.WebView) {
    bind := func(name string, fn interface{}) {
        if err := w.Bind(name, fn); err != nil {
            fmt.Printf("Failed to bind %s: %v\n", name, err)
        }
    }

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

    type aTkn struct {
        Name string `json:"name"`
        Uuid string `json:"uuid"`
        Xuid string `json:"xuid"`
        ClientID string `json:"clientID"`
        AccessToken string `json:"accessToken"`
        UserType string `json:"userType"`
    }

    bind("WriteFile", writeFile)
    bind("LocalAppdata", func() string { return LocalAppData })
    bind("RoamingAppdata", func() string { return RoamingAppData })
    bind("edition", func() string { return edition })
    bind("login_with_microsft", func() { 
        w.Dispatch(func() {
            w.SetTitle("Minecraft Login")
            w.Navigate("https://sisu.xboxlive.com/connect/XboxLive/?state=login&cobrandId=8058f65d-ce06-4c30-9559-473c9275a65d&tid=896928775&ru=https%3A%2F%2Fwww.minecraft.net%2Fen-us%2Flogin&aid=1142970254")
            w.Bind("println", func (fn any)  {
                fmt.Println(fn)
            })
            w.Bind("save", func (data string)  {
                os.WriteFile("cookie.txt", []byte(data), fs.FileMode(os.O_CREATE))
                q, _ := regexp.Compile("token:.*:")
                AccessToken := strings.ReplaceAll(strings.ReplaceAll(q.FindString(data), "token:", ""), ":", "")
                jwtData, _ := base64.RawStdEncoding.DecodeString(strings.Split(AccessToken, ".")[1])
                Uuid, _ := jsonparser.GetString(jwtData, "profiles", "mc")
                resp, _ := http.Get("https://playerdb.co/api/player/minecraft/"+Uuid)
                defer resp.Body.Close()
                body := make([]byte, 0)
                if (resp.StatusCode == http.StatusOK) { body, _ = io.ReadAll(resp.Body) }
                Name, _ := jsonparser.GetString(body, "data", "player", "username")
                Xuid, _ := jsonparser.GetString(jwtData, "xuid")
                ClientID := base64.RawStdEncoding.EncodeToString([]byte(uuid.NewV1().String()))
                // YmQyNDViY2UtOGFhNC00ODNmLWI4NzctNmFiZmIxZWE4MWY5
                
                UserData := aTkn{
                    Name: Name,
                    Uuid: Uuid,
                    Xuid: Xuid,
                    ClientID: ClientID,
                    AccessToken: AccessToken,
                    UserType: "msa",
                }

                bytez, _ := json.MarshalIndent(UserData, "", "  ")
                os.WriteFile(RoamingAppData + "/.crackerClient/AuthConfig.json", bytez, os.FileMode(os.O_CREATE))

            })
            go func() {
                stop := false
                w.Dispatch(func() {
                    w.Bind("back", func() {
                        stop = true
                        w.SetTitle("Cracker Client")
                        w.Navigate("http://localhost:"+strconv.Itoa(port))
                    })
                })
                for {
                    if (stop) {
                        break
                    }
                    fmt.Println(stop)
                    w.Dispatch(func() {
                        w.Eval("setTimeout(()=>{if (document.domain == 'www.minecraft.net') { save(document.cookie); back() }}, 6000)")
                    })
                    time.Sleep(1 * time.Second)
                }
            }()
        })
     })
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
