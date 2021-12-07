package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	config := &Config{}

	if data, err := os.ReadFile("config.json"); err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
	}

	client := twitch.NewAnonymousClient()
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		text := message.Message
		words := strings.Split(text, " ")
		for i, v := range words {
			link, ok := config.BetterTTVEmotes[v]
			if ok {
				words[i] = `<img src="` + link + `"  width="25" height="25">`
			}
		}
		text = strings.Join(words, " ")
		for _, v := range message.Emotes {
			text = strings.Replace(text, v.Name, fmt.Sprintf(`<img src="https://static-cdn.jtvnw.net/emoticons/v2/%s/static/dark/1.0" width="25" height="25">`, v.ID), v.Count)
		}
		message.Message = text
		message.Channel = `<img src="` + config.AvatarMap[strings.ToLower(message.Channel)] + `" width="25" height="25">`
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

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/chat.html")
}
