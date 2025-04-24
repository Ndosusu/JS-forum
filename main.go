package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/RealTimeForumDosJean/server/api"
	localDatabase "github.com/RealTimeForumDosJean/server/database"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func serveFile(filename string, w http.ResponseWriter, data any, fatal bool) {
	index, errTmpl := template.ParseFiles(filename)
	if errTmpl != nil {
		log.Printf("Error creating template for %s: %v\n", filename, errTmpl)

		if fatal {
			os.Exit(1)
		}
	}

	if errExec := index.Execute(w, data); errExec != nil {
		log.Printf("Error executing template for %s: %v\n", filename, errExec)

		if fatal {
			os.Exit(1)
		}
	}
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	serveFile("index.html", w, nil, true)
}

func main() {
	localDatabase.InitConnection()

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	// ------------------------ HANDLERS HOOKING ------------------------
	mux.HandleFunc("/", indexHandler)

	// LOGIN / REGISTER
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/imConnected", api.LoginFromCookieHandler)
	mux.HandleFunc("/api/register", api.RegisterHandler)
	mux.HandleFunc("/api/logout", api.LogoutHandler)

	// USER-RELATED
	mux.HandleFunc("/api/fetchUsers", api.FetchAllUsersHandler)
	mux.HandleFunc("/api/whosConnected", api.FetchConnectedUsers)

	// POST-RELATED
	mux.HandleFunc("/api/newPost", api.CreatePostHandler)
	mux.HandleFunc("/api/fetchAllPosts", api.FetchAllPostsHandler)
	mux.HandleFunc("/api/postFetcher", api.FetchPostHandler)

	// COMMENT-RELATED
	mux.HandleFunc("/api/newComment", api.CreateCommentHandler)
	mux.HandleFunc("/api/fetchComments", api.FetchAllCommentsHandler)
	mux.HandleFunc("/api/commentFetcher", api.FetchCommentHandler)

	mux.HandleFunc("/ws", handleWS)
	// ------------------------------------------------------------------

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
