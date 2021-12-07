package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"twitch-multichat/utils"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	config := &Config{}

	if data, err := os.ReadFile("config.json"); err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
	}
	tw := utils.TwitchClient{ID: config.TwitchID, Token: config.TwitchSecret}
	data, _ := tw.GrabUserInfo(http.DefaultClient, "DepressedHumanist")
	if data == nil {
		http.Error(w, "Nope", 500)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(""))
	}
}

func Token(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/token.html")
		return
	}
	twitchBullshit := make([]byte, 1024*10)
	n, err := r.Body.Read(twitchBullshit)
	twitchBullshitData := string(twitchBullshit[:n])
	accessToken := ""
	for _, v := range strings.Split(twitchBullshitData, "&") {
		kv := strings.Split(v, "=")
		if kv[0] == "access_token" {
			accessToken = kv[1]
		}
	}

	config := &Config{}

	if data, err := os.ReadFile("config.json"); err == nil {
		err = json.Unmarshal(data, &config)
		if err != nil {
			http.Error(w, "bad request", 400)
			return
		}
	}
	config.TwitchID = "yig9b3sexzjwm2wq2y9ap6c6lp5o80"
	config.TwitchSecret = accessToken

	data, _ := json.Marshal(config)
	err = os.WriteFile("config.json", data, os.ModePerm)
	if err != nil {
		http.Error(w, "can't write config", 500)
		return
	}
	w.WriteHeader(200)
	_, _ = fmt.Fprint(w, "Done")
}
