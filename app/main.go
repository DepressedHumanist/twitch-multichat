package main

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"net/http"
	"os"
	"twitch-multichat/handlers"
	"twitch-multichat/utils"
)

func main() {
	http.HandleFunc("/chat", handlers.ServeHome)
	http.HandleFunc("/home", handlers.Index)
	http.HandleFunc("/token", handlers.Token)
	http.HandleFunc("/token_check", handlers.Ping)
	http.HandleFunc("/ws", handlers.WsHandler)
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	go http.ListenAndServe(":8081", nil)
	var a, _ = astilectron.New(nil, astilectron.Options{
		AppName:            "multichat",
		BaseDirectoryPath:  "deps",
		VersionAstilectron: "0.33.0",
		VersionElectron:    "6.1.2",
	})
	defer a.Close()

	// Start astilectron
	a.Start()
	useContentSize := true
	url := "http://localhost:8081/home"
	if _, err := os.Stat("config.json"); err != nil {
		url = "http://localhost:8081/home"
		utils.OpenBrowser(utils.TOKEN_REFRESH_URL)
	}

	var w, _ = a.NewWindow(url, &astilectron.WindowOptions{
		UseContentSize: &useContentSize,
		Height:         astikit.IntPtr(600),
		Width:          astikit.IntPtr(400),
	})
	w.Create()

	// Blocking pattern
	a.Wait()

}
