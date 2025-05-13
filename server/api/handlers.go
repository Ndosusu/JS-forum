package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	forumServer "github.com/RealTimeForumDosJean/server"

	dbComment "github.com/RealTimeForumDosJean/server/structs/comment"
	dbDM "github.com/RealTimeForumDosJean/server/structs/directMsg"
	dbPost "github.com/RealTimeForumDosJean/server/structs/post"
	dbUser "github.com/RealTimeForumDosJean/server/structs/user"
)

// ----------------------- CONNECTED USERS HANDLER -----------------------

type ClientsHandler struct {
	ConnectedClients  map[string]*http.Cookie
	ModificationMutex sync.Mutex
}

var CliHandler = ClientsHandler{
	ConnectedClients:  make(map[string]*http.Cookie, 0),
	ModificationMutex: sync.Mutex{},
}

func (cliHandler *ClientsHandler) AddUser(newUsr dbUser.User, usrCookie *http.Cookie) bool {
	couldLock := cliHandler.ModificationMutex.TryLock()
	if !couldLock {
		return couldLock
	}

	cliHandler.ConnectedClients[newUsr.UUID] = usrCookie

	cliHandler.ModificationMutex.Unlock()
	return couldLock
}

func (cliHandler *ClientsHandler) RemoveUser(userUUID string) bool {
	couldLock := cliHandler.ModificationMutex.TryLock()
	if !couldLock {
		return couldLock
	}

	for uuid := range cliHandler.ConnectedClients {
		if uuid == userUUID {
			delete(cliHandler.ConnectedClients, uuid)
			break
		}
	}

	cliHandler.ModificationMutex.Unlock()
	return couldLock
}

func (cliHandler *ClientsHandler) IsUserConnected(uuid string) bool {
	for userUUID := range cliHandler.ConnectedClients {
		if userUUID == uuid {
			errValid := cliHandler.ConnectedClients[userUUID].Valid()
			if errValid != nil {
				cliHandler.RemoveUser(userUUID)
			}
			return errValid == nil
		}
	}

	return false
}

func LoginFromCookieHandler(w http.ResponseWriter, r *http.Request) {
	usrCookie, errCookie := r.Cookie("UserLogged")
	if errCookie != nil {
		http.Error(w, "No cookie", http.StatusForbidden)
		return
	}

	for _, cookie := range CliHandler.ConnectedClients {
		if usrCookie.Value == cookie.Value {
			http.Error(w, "You're already connected for me", http.StatusOK)
			return
		}
	}

	parts := strings.Split(usrCookie.Value, "|")

	usr, errFetch := dbUser.FetchUserByName(r.Context(), parts[1])
	if errFetch != nil {
		http.Error(w, "Error while fetching, please log in again", http.StatusInternalServerError)
		return
	}

	if usr == (dbUser.User{}) {
		http.Error(w, "I don't know you", http.StatusForbidden)
		return
	}
	if usr.UUID == parts[0] {
		CliHandler.AddUser(usr, usrCookie)
		http.Error(w, "You're right", http.StatusOK)
	}
}

// -----------------------------------------------------------------------

// ------------------------ USER-RELATED HANDLERS ------------------------

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "{\"Error\": \"Method not allowed\"}", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	var reqBody []byte
	var usr dbUser.User
	var err error

	if reqBody, err = io.ReadAll(r.Body); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error body\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if len(reqBody) == 0 {
		http.Error(w, "{\"Error\": \"Body empty\"}", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, "Body is empty")
		return
	}

	if err = json.Unmarshal(reqBody, &usr); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error marshal\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error unmarshaling: "+err.Error())
		return
	}

	{
		var usrFound dbUser.User

		// Does user exist?
		usrFound, err = dbUser.FetchUserByEmail(r.Context(), usr.Email)
		if err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "Error fetching email: "+err.Error())
			return
		}

		// No => return error
		if usrFound == (dbUser.User{}) {
			http.Error(w, "{\"Error\": \"No user was found with this email address. Please try another.\"}", http.StatusNotFound)
			fmt.Fprintln(os.Stderr, "No user found")
			return
		}

		// Yes but password invalid => return error
		if err = dbUser.CheckPassword(usrFound.EncryptedPassword, usr.EncryptedPassword); err != nil {
			http.Error(w, "{\"Error\": \"Password did not match. Please try another.\"}", http.StatusUnauthorized)
			fmt.Fprintf(os.Stderr, "Invalid password for user %s\n", usrFound.Username)
			return
		}

		usr = usrFound
	}

	fmt.Printf("User logged in: %s -> %s (%s)\n", usr.UUID, usr.Username, usr.Email)

	newCookie := http.Cookie{
		Name:   "UserLogged",
		Path:   "/",
		Value:  usr.ToCookieValue(),
		MaxAge: 600, // 10 minutes
	}

	http.SetCookie(w, &newCookie)
	CliHandler.AddUser(usr, &newCookie)

	w.WriteHeader(http.StatusOK)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "{\"Error\": \"Method not allowed\"}", http.StatusBadRequest)
		fmt.Fprintln(os.Stderr, r.Form)
		return
	}

	var newUser dbUser.User
	var err error

	{
		var reqBody []byte

		if reqBody, err = io.ReadAll(r.Body); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error body\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "Error fetching body: "+err.Error())
			return
		}

		if len(reqBody) == 0 {
			http.Error(w, "{\"Error\": \"Body empty\"}", http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(reqBody, &newUser); err != nil {
			http.Error(w, "{\"Error\": \"Fatal error unmarshal\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "Error unmarshaling: "+err.Error())
			return
		}
	}

	// Verify that user doesn't exist already
	{
		var usrFound dbUser.User
		var exists bool

		usrFound, err = dbUser.FetchUserByEmail(r.Context(), newUser.Email)
		if err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "Error fetching email: "+err.Error())
			return
		}

		if usrFound != (dbUser.User{}) {
			http.Error(w, "{\"Error\": \"User already exists with this email address. Please try with another.\"}", http.StatusUnauthorized)
			return
		}

		exists, err = dbUser.IsUsernameTaken(r.Context(), newUser.Username)
		if err != nil {
			http.Error(w, "{\"Error\": \"Fatal error fetching\"}", http.StatusInternalServerError)
			fmt.Fprintln(os.Stderr, "Error fetching username: "+err.Error())
			return
		}

		if exists {
			http.Error(w, "{\"Error\": \"User already exists with this username. Please try with another.\"}", http.StatusUnauthorized)
			return
		}
	}

	var uuid string

	if uuid, err = forumServer.GenerateUUID(); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error gen\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error generating ID: "+err.Error())
		return
	}

	if newUser.EncryptedPassword, err = dbUser.HashPassword(newUser.EncryptedPassword); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error hash\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error hashing password: "+err.Error())
		return
	}

	newUser.UUID = uuid
	newUser.CreatedAt = time.Now()
	newUser.Role = "user"

	if err = dbUser.RegisterUser(r.Context(), newUser.ToMap()); err != nil {
		http.Error(w, "{\"Error\": \"Fatal error add\"}", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, "Error adding to database: "+err.Error())
		return
	}

	fmt.Printf("New user registered: %s -> %s (%s)\n", newUser.UUID, newUser.Username, newUser.Email)

	newCookie := http.Cookie{
		Name:   "UserLogged",
		Path:   "/",
		Value:  newUser.ToCookieValue(),
		MaxAge: 600, // 10 minutes
	}

	http.SetCookie(w, &newCookie)
	CliHandler.AddUser(newUser, &newCookie)

	w.WriteHeader(http.StatusOK)
}

func PP_Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usr dbUser.User

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		http.Error(w, "Fatal error decode id", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture, err = dbUser.FetchPPByID(r.Context(), usr.UUID); err != nil {
		http.Error(w, "Fatal error query pp", http.StatusInternalServerError)
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if usr.ProfilePicture == "" {
		fmt.Printf("User not found for ID \"%s\"\n", usr.UUID)
	}

	json.NewEncoder(w).Encode(usr.ProfilePicture)
}

func LogMeOut(r *http.Request) error {
	uuid, errCookie := forumServer.GetUserUUIDFromCookie(r)
	if errCookie != nil {
		return fmt.Errorf("error getting cookie: %s", errCookie.Error())
	}

	fmt.Printf("User with ID %s requested logout\n", uuid)
	CliHandler.RemoveUser(uuid)

	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if errLogOut := LogMeOut(r); errLogOut != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:    "UserLogged",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func FetchAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer tous les utilisateurs
	userData, err := dbUser.FetchAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

func FetchConnectedUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	offlineUsers := make([]int, 0)

	// Récupérer tous les utilisateurs
	connectedUsers, err := dbUser.FetchAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, user := range connectedUsers {
		match := false
		for uuid := range CliHandler.ConnectedClients {
			if user.UUID == uuid {
				match = true
				break
			}
		}

		if !match {
			offlineUsers = append(offlineUsers, i)
		}
	}

	for len(offlineUsers) > 0 {
		connectedUsers = slices.Delete(connectedUsers, offlineUsers[len(offlineUsers)-1], offlineUsers[len(offlineUsers)-1]+1)
		offlineUsers = offlineUsers[:len(offlineUsers)-1]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connectedUsers)
}

// -----------------------------------------------------------------------

// --------------------- HANDLERS FOR DIRECT MESSAGE ---------------------

func FetchDMsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	usrCookie, cookieErr := r.Cookie("UserLogged")
	if cookieErr != nil {
		http.Error(w, "No cookie detected", http.StatusForbidden)
		return
	}

	idRegex, regexErr := regexp.Compile("[0-9a-z-]+")
	if regexErr != nil {
		http.Error(w, regexErr.Error(), http.StatusInternalServerError)
		return
	}

	match := idRegex.Find([]byte(usrCookie.Value))
	if match == nil {
		http.Error(w, "Invalid UUID in cookie", http.StatusBadRequest)
		return
	}

	myDMs, errFetch := dbDM.FetchMyDMs(r.Context(), string(match))
	if errFetch != nil {
		http.Error(w, errFetch.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(myDMs)
}

// ------------------------ POST-RELATED HANDLERS ------------------------

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newPost dbPost.Post
	newPost.Likes = 0
	newPost.Dislikes = newPost.Likes
	newPost.Created_at = time.Now()

	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	post_UUID, errGen := forumServer.GenerateUUID()
	if errGen != nil {
		http.Error(w, "erreur lors de la génération de l'UUID du post: "+errGen.Error(), http.StatusInternalServerError)
		return
	}

	newPost.UUID = post_UUID

	// extraire le uuid du cookie
	user_UUID, err := forumServer.GetUserUUIDFromCookie(r)
	if err != nil {
		http.Error(w, "erreur lors de la récupération du uuid: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newPost.User_uuid = user_UUID

	if newPost.Title == "" || newPost.Content == "" {
		http.Error(w, "informations manquantes", http.StatusBadRequest)
		return
	}

	err = dbPost.CreatePost(r, newPost.ToMap())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(w, "Post created", http.StatusOK)
}

func FetchPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	userUUID := queryParams.Get("user_uuid")
	postUUID := queryParams.Get("post_uuid")

	params := map[string]any{}

	if userUUID != "" {
		params["user_uuid"] = userUUID
	} else if postUUID != "" {
		params["post_uuid"] = postUUID
	}

	postData, err := dbPost.FetchPost(r.Context(), params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}

func FetchAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer tous les posts
	postData, err := dbPost.FetchAllPosts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postData)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer le post_uuid depuis les paramètres de la requête
	var params map[string]any

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbPost.DeletePost(r.Context(), params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner une réponse succès
	w.WriteHeader(http.StatusNoContent) // No Content, car aucune donnée à renvoyer
}

// -----------------------------------------------------------------------

// ---------------------- COMMENTS-RELATED HANDLERS ----------------------

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var comment struct {
		PostUUID string `json:"post_uuid"`
		Content  string `json:"content"`
		UserUUID string `json:"user_uuid"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if comment.PostUUID == "" || comment.Content == "" || comment.UserUUID == "" || comment.Username == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	createdComment, err := dbComment.CreateComment(r.Context(), map[string]any{
		"post_uuid": comment.PostUUID,
		"content":   comment.Content,
		"user_uuid": comment.UserUUID,
		"username":  comment.Username,
	})
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdComment)
}

func FetchCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params map[string]any

	// Récupere la demande du front et decode le JSON post_uuid
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Vérification de la présence de post_uuid
	postUUID, ok := params["post_uuid"].(string)
	if !ok || postUUID == "" {
		http.Error(w, "Missing or invalid post_uuid", http.StatusBadRequest)
		return
	}

	// Récupération des commentaires basés sur le post_uuid
	commentData, err := dbComment.FetchComment(r.Context(), map[string]any{
		"post_uuid": postUUID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si trouvé renvoie la réponse en format JSON pour le fronted etc....
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(commentData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func FetchCommentsByPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postUUID := r.URL.Query().Get("post_uuid")
	if postUUID == "" {
		http.Error(w, "Missing post_uuid parameter", http.StatusBadRequest)
		return
	}

	comments, err := dbComment.FetchComment(r.Context(), map[string]any{"post_uuid": postUUID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func FetchAllCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer tous les posts
	commentData, err := dbComment.FetchAllComments(r.Context())

	//log.Println(commentData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commentData)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer le post_uuid depuis les paramètres de la requête
	var params map[string]any

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbComment.DeleteComment(r.Context(), params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner une réponse succès
	w.WriteHeader(http.StatusNoContent) // No Content, car aucune donnée à renvoyer
}

// -----------------------------------------------------------------------
