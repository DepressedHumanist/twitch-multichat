package main

import (
	"github.com/lxn/win"
	"github.com/webview/webview"
	"net/http"
	"os"
	"runtime"
	"twitch-multichat/handlers"
)

func main() {
	http.HandleFunc("/chat", handlers.ServeHome)
	http.HandleFunc("/home", handlers.Index)
	http.HandleFunc("/token", handlers.Token)
	http.HandleFunc("/token_check", handlers.Ping)
	http.HandleFunc("/ws", handlers.WsHandler)
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	go http.ListenAndServe(":8081", nil)
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Multichat")
	w.SetSize(400, 600, webview.HintNone)
	if _, err := os.Stat("config.json"); err == nil {
		w.Navigate("http://localhost:8081/home")
	} else {
		w.Navigate("http://localhost:8081/token")
	}

	if runtime.GOOS == "windows" {
		winflags := win.GetWindowLong(win.HWND(w.Window()), -20)
		winflags &= ^(win.WS_CAPTION | win.WS_THICKFRAME | win.WS_MINIMIZEBOX | win.WS_MAXIMIZEBOX | win.WS_SYSMENU)
		win.SetWindowLong(win.HWND(w.Window()), -20, winflags)
	}
	w.Run()
}
