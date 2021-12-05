package main

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/websocket"
	"github.com/lxn/win"
	"github.com/webview/webview"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Config struct {
	TwitchID     string   `json:"twitchID"`
	TwitchSecret string   `json:"twitchSecret"`
	Streamers    []string `json:"streamers"`
}

func handler(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadFile("secrets.json")
	// can file be opened?
	if err != nil {
		fmt.Print(err)
	}
	config := &Config{}
	_ = json.Unmarshal(b, config)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	twi := NewTwitchClient(config.TwitchID, config.TwitchSecret)
	chanIds := make([]string, len(config.Streamers))
	avatarMap := make(map[string]string)
	for i, v := range config.Streamers {
		user := twi.grabUserInfo(http.DefaultClient, v)
		//todo handle channel not found
		avatarMap[strings.ToLower(user.Name)] = user.ProfilePic
		chanIds[i] = user.ID
	}
	bttvEmotes := GetBetterTTVEmotes(chanIds...)
	log.Println(bttvEmotes)
	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		text := message.Message
		words := strings.Split(text, " ")
		for i, v := range words {
			link, ok := bttvEmotes[v]
			if ok {
				words[i] = `<img src="` + link + `"  width="25" height="25">`
			}
		}
		text = strings.Join(words, " ")
		for _, v := range message.Emotes {
			text = strings.Replace(text, v.Name, fmt.Sprintf(`<img src="https://static-cdn.jtvnw.net/emoticons/v2/%s/static/dark/1.0" width="25" height="25">`, v.ID), v.Count)
		}
		message.Message = text
		message.Channel = `<img src="` + avatarMap[strings.ToLower(message.Channel)] + `" width="25" height="25">`
		err = conn.WriteJSON(message)
		if err != nil {
			log.Println(err.Error())
		}
	})
	client.Join(config.Streamers...)

	if err = client.Connect(); err != nil {
		panic(err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/chat.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handler)
	http.Handle("/templates", http.FileServer(http.Dir("templates")))
	go http.ListenAndServe(":8081", nil)
	debug := false
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Multichat")
	w.SetSize(800, 400, webview.HintNone)
	w.Navigate("http://localhost:8081")
	winflags := win.GetWindowLong(win.HWND(w.Window()), -20)
	winflags &= ^(win.WS_CAPTION | win.WS_THICKFRAME | win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX | win.WS_SYSMENU)
	win.SetWindowLong(win.HWND(w.Window()), -20, winflags)
	w.Run()
}
