package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"twitch-multichat/utils"
)

type Config struct {
	TwitchID        string            `json:"twitchID"`
	TwitchSecret    string            `json:"twitchSecret"`
	Streamers       []string          `json:"streamers"`
	SkipInit        bool              `json:"skipInit"`
	AvatarMap       map[string]string `json:"avartarMap"`
	BetterTTVEmotes map[string]string `json:"betterTTVEmotes"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	config := &Config{}

	if data, err := os.ReadFile("config.json"); err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
	}
	if r.Method == "GET" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		if _, ok := r.Form["force"]; !ok && config.SkipInit {
			http.Redirect(w, r, "/chat", 302)
		}
		files, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "can't read template "+err.Error(), 500)
			return
		}
		err = files.Execute(w, config)
		if err != nil {
			http.Error(w, "can't execute template", 500)
			return
		}
		return
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
		if streamers, ok := r.Form["streamer"]; ok {
			config.Streamers = streamers
		}
		_, ok := r.Form["skipInit"]
		config.SkipInit = ok

		twi := utils.TwitchClient{ID: config.TwitchID, Token: config.TwitchSecret}
		chanIds := make([]string, len(config.Streamers))
		config.AvatarMap = make(map[string]string)
		for i, v := range config.Streamers {
			user, err := twi.GrabUserInfo(http.DefaultClient, v)
			if err != nil {
				data, _ := json.Marshal(config)
				err = os.WriteFile("config.json", data, os.ModePerm)
				if err != nil {
					http.Error(w, "can't write config", 500)
					return
				}
				utils.OpenBrowser(utils.TOKEN_REFRESH_URL)
				http.Redirect(w, r, "/token", 302)
				return
			}
			//todo handle channel not found
			config.AvatarMap[strings.ToLower(user.Name)] = user.ProfilePic
			chanIds[i] = user.ID
		}
		config.BetterTTVEmotes = utils.GetBetterTTVEmotes(chanIds...)

		data, _ := json.Marshal(config)
		err = os.WriteFile("config.json", data, os.ModePerm)
		if err != nil {
			http.Error(w, "can't write config", 500)
			return
		}
		http.Redirect(w, r, "/chat", 302)
		return
	}
}
