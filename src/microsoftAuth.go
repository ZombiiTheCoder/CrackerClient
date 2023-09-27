package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/twinj/uuid"
	webview "github.com/webview/webview_go"
)

type AuthStruct struct {
	Name        string `json:"name"`
	Uuid        string `json:"uuid"`
	Xuid        string `json:"xuid"`
	ClientID    string `json:"clientID"`
	AccessToken string `json:"accessToken"`
	UserType    string `json:"userType"`
}

var login_with_microsoft = func(w webview.WebView) {
	w.Dispatch(func() {
		w.SetTitle("Minecraft Login")
		w.Navigate("https://sisu.xboxlive.com/connect/XboxLive/?state=login&cobrandId=8058f65d-ce06-4c30-9559-473c9275a65d&tid=896928775&ru=https%3A%2F%2Fwww.minecraft.net%2Fen-us%2Flogin&aid=1142970254")
		w.Bind("println", func(fn any) {
			fmt.Println(fn)
		})
		w.Bind("save", func(data string) {
			q, _ := regexp.Compile("token:.*:")
			AccessToken := strings.ReplaceAll(strings.ReplaceAll(q.FindString(data), "token:", ""), ":", "")
			jwtData, _ := base64.RawStdEncoding.DecodeString(strings.Split(AccessToken, ".")[1])
			Uuid, _ := jsonparser.GetString(jwtData, "profiles", "mc")
			resp, _ := http.Get("https://playerdb.co/api/player/minecraft/" + Uuid)
			defer resp.Body.Close()
			body := make([]byte, 0)
			if resp.StatusCode == http.StatusOK {
				body, _ = io.ReadAll(resp.Body)
			}
			Name, _ := jsonparser.GetString(body, "data", "player", "username")
			Xuid, _ := jsonparser.GetString(jwtData, "xuid")
			ClientID := base64.RawStdEncoding.EncodeToString([]byte(uuid.NewV1().String()))

			UserData := AuthStruct{
				Name:        Name,
				Uuid:        Uuid,
				Xuid:        Xuid,
				ClientID:    ClientID,
				AccessToken: AccessToken,
				UserType:    "msa",
			}

			bytez, _ := json.MarshalIndent(UserData, "", "  ")
			os.WriteFile(RoamingAppData+"/.crackerClient/AuthConfig.json", bytez, os.FileMode(os.O_CREATE))

		})
		go func() {
			stop := false
			w.Dispatch(func() {
				w.Bind("back", func() {
					stop = true
					w.SetTitle("Cracker Client")
					w.Navigate("http://localhost:" + strconv.Itoa(port))
				})
			})
			for {
				if stop {
					break
				}
				fmt.Println(stop)
				w.Dispatch(func() {
					w.Eval("setTimeout(()=>{if (location.href == 'https://www.minecraft.net/en-us/msaprofile') { save(document.cookie); back() }}, 2000)")
				})
				time.Sleep(1 * time.Second)
			}
		}()
	})
}