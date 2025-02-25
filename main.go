package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

// Configuration WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erreur WebSocket :", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connexion WebSocket fermée :", err)
			break
		}
		log.Printf("Message reçu : %s", msg)
		conn.WriteMessage(websocket.TextMessage, []byte("Message reçu !"))
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handleWS)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
