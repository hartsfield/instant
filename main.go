package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var dat []byte

func main() {
	dat, _ = os.ReadFile("words")
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			words := strings.Join(getWords(msg, dat), ", ")

			// Write message back to browser
			if err = conn.WriteMessage(msgType, []byte(words)); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.tmpl")
	})

	http.ListenAndServe(":8085", nil)
}

func getWords(msg, data []byte) []string {
	words := strings.Split(string(data), "\n")
	var beg, end int
	unfound := true
	for i, word := range words {
		if len(word) >= len(string(msg)) {
			if strings.ToLower(word[0:len(string(msg))]) == strings.ToLower(string(msg)) && unfound {
				unfound = false
				beg = i
			}
			if strings.ToLower(word[0:len(string(msg))]) != strings.ToLower(string(msg)) && beg > 0 {
				end = i
				return words[beg:end]
			}
		}
	}
	return words[beg:end]
}
