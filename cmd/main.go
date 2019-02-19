package main

import (
	"log"
	"net/http"
	"text/template"

	chat "github.com/brunobandev/face-to-face/chat"
)

var index = template.Must(template.ParseFiles("./index.html"))

func home(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

func main() {
	go chat.DefaultHub.Start()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", chat.WSHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
