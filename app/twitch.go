package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// JSON handling

type TwitchUsers struct {
	Data []TwitchUser `json:"data"`
}

type TwitchUser struct {
	ID         string `json:"id"`
	Name       string `json:"display_name"`
	ProfilePic string `json:"profile_image_url"`
}

type BetterTTVUser struct {
	ChannelEmotes []BetterTTVEmote `json:"channelEmotes"`
	SharedEmotes  []BetterTTVEmote `json:"sharedEmotes"`
}

type BetterTTVEmote struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type FrankerFaceZEmote struct {
	Code   string   `json:"code"`
	Images FFZImage `json:"images"`
}

type FFZImage struct {
	X1 string `json:"1x"`
}

// TwitchClient Twitch api client
type TwitchClient struct {
	ID    string
	Token string `json:"access_token"`
}

func (t TwitchClient) grabUserInfo(client *http.Client, name string) *TwitchUser {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login="+name, nil)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	req.Header.Add("Client-Id", t.ID)
	req.Header.Add("Authorization", "Bearer "+t.Token)
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error: couldn't get a user from twitch: ", err.Error())
		return nil
	}
	data := make([]byte, 1024)
	n, err := res.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Fatal(err.Error())
	}
	userData := &TwitchUsers{}
	err = json.Unmarshal(data[:n], userData)
	if err != nil {
		log.Println("Error: got broken json from twitch: ", err.Error())
		return nil
	}
	if len(userData.Data) < 1 {
		log.Println("Error: user not found")
		return nil
	}
	return &userData.Data[0]
}

func GetBetterTTVEmotes(channel ...string) map[string]string {
	res, err := http.Get("https://api.betterttv.net/3/cached/emotes/global")
	if err != nil {
		return nil
	}
	data := make([]byte, 100*1024)
	n, err := res.Body.Read(data)
	if err != nil && err != io.EOF {
		log.Fatal(err.Error())
	}
	emotes := make([]*BetterTTVEmote, 0)
	err = json.Unmarshal(data[:n], &emotes)
	resMap := make(map[string]string)
	for _, v := range emotes {
		resMap[v.Code] = "https://cdn.betterttv.net/emote/" + v.ID + "/1x"
	}
	for _, v := range channel {
		// BTTV
		res, err = http.Get("https://api.betterttv.net/3/cached/users/twitch/" + v)
		if err != nil {
			return nil
		}
		data = make([]byte, 100*1024)
		n, err = res.Body.Read(data)
		if err != nil && err != io.EOF {
			log.Fatal(err.Error())
		}
		user := &BetterTTVUser{}
		err = json.Unmarshal(data[:n], user)
		for _, emote := range user.ChannelEmotes {
			resMap[emote.Code] = "https://cdn.betterttv.net/emote/" + emote.ID + "/1x"
		}
		for _, emote := range user.ChannelEmotes {
			resMap[emote.Code] = "https://cdn.betterttv.net/emote/" + emote.ID + "/1x"
		}

		// FFZ
		res, err = http.Get("https://api.betterttv.net/3/cached/frankerfacez/users/twitch/" + v)
		if err != nil {
			return nil
		}
		data = make([]byte, 100*1024)
		n, err = res.Body.Read(data)
		if err != nil && err != io.EOF {
			log.Fatal(err.Error())
		}
		ffzImages := make([]*FrankerFaceZEmote, 0)
		err = json.Unmarshal(data[:n], &ffzImages)
		for _, emote := range ffzImages {
			resMap[emote.Code] = emote.Images.X1
		}
	}
	return resMap
}

func NewTwitchClient(cliId, secret string) *TwitchClient {
	// TODO read from a file or something
	res, err := http.Post(fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", cliId, secret), "", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	cli := &TwitchClient{
		ID: cliId,
	}
	data := make([]byte, 1024)
	n, err := res.Body.Read(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = json.Unmarshal(data[:n], cli)
	if err != nil {
		log.Fatal(err.Error())
	}
	return cli
}
