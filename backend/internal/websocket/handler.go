package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler gerencia conexoes WebSocket para assistente em tempo real
func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("cliente websocket conectado")

	// Loop de leitura do cliente
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("websocket read error: %v", err)
			return
		}
	}
}
