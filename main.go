package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fmt.Println("Socket running in localhost:8080/socket")

	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}
			input := string(message)

			content, err := ioutil.ReadFile("./files/" + input)
			if err != nil {
				err = conn.WriteMessage(mt, []byte("404 - Not Found"))
				if err != nil {
					log.Println("write failed:", err)
				}
				continue
			}

			output := "200 - OK"
			output += "\n----------------------------------------"
			output += "\n" + string(content)
			output += "\n----------------------------------------"
			err = conn.WriteMessage(mt, []byte(output))
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("port 8080 already used")
	}
}
