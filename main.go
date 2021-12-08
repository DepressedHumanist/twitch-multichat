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
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", handlers.Chat)
	mux.HandleFunc("/token", handlers.Token)
	mux.HandleFunc("/token_check", handlers.Ping)
	mux.HandleFunc("/ws", handlers.WsHandler)
	mux.HandleFunc("/", handlers.Index)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	go http.ListenAndServe(":8081", mux)
	var a, _ = astilectron.New(nil, astilectron.Options{
		AppName:            "multichat",
		AppIconDefaultPath: "static/icon.png",
		AppIconDarwinPath:  "static/icon.icns",
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
		url = "http://localhost:8081/token"
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
