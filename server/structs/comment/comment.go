package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	forumServer "github.com/RealTimeForumDosJean/server"
	localDatabase "github.com/RealTimeForumDosJean/server/database"
)

type Comment struct {
	Comment_id     string    `json:"comment_id"`
	Post_uuid      string    `json:"post_uuid"`
	User_uuid      string    `json:"user_uuid"`
	Content        string    `json:"content"`
	Username       string    `json:"username"`
	ProfilePicture string    `json:"profile_picture"`
	Created_at     time.Time `json:"created_at"`
}

func fromMap(values map[string]any) Comment {
	newComment := Comment{}

	if commentUUID, ok := values["comment_id"].(string); ok {
		newComment.Comment_id = commentUUID
	}
	if postUUID, ok := values["post_uuid"].(string); ok {
		newComment.Post_uuid = postUUID
	}
	if userUUID, ok := values["user_uuid"].(string); ok {
		newComment.User_uuid = userUUID
	}
	if content, ok := values["content"].(string); ok {
		newComment.Content = content
	}
	if createdAt, ok := values["created_at"].(time.Time); ok {
		newComment.Created_at = createdAt
	}
	if username, ok := values["username"].(string); ok {
		newComment.Username = username
	}
	if profilePicture, ok := values["profile_picture"].(string); ok {
		newComment.ProfilePicture = profilePicture
	}

	return newComment
}

// CreateComment crée un nouveau commentaire et l'insère dans la base de données
func CreateComment(ctx context.Context, params map[string]any) (*Comment, error) {

	comment_UUID, _ := forumServer.GenerateUUID()

	fmt.Println(comment_UUID)

	post_UUID, postOK := params["post_uuid"].(string)

	// !userOK
	user_UUID, _ := params["user_uuid"].(string)
	content, contentOK := params["content"].(string)

	//!userOK

	if !postOK || !contentOK {
		return nil, errors.New("informations manquantes pour le commentaire")
	}

	creationDate := time.Now()

	// Insertion du commentaire dans la table
	createCommentQuery := `INSERT INTO Comments (comment_id, post_uuid, user_uuid, content, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := localDatabase.RunDatabaseQuery(ctx, createCommentQuery, comment_UUID, post_UUID, user_UUID, content, creationDate, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du commentaire: %v", err)
	}

	// Création de la structure Comment
	newComment := &Comment{
		Comment_id: comment_UUID,
		Post_uuid:  post_UUID,
		User_uuid:  user_UUID,
		Content:    content,
		Created_at: creationDate,
	}

	return newComment, nil
}

// FetchAllComments récupère tous les commentaires de la base de données
func FetchAllComments(ctx context.Context) ([]Comment, error) {
	fetchAllCommentsQuery := `
		SELECT c.*, u.username, u.profile_picture
		FROM Comments c
		JOIN Users u ON c.user_uuid = u.user_uuid
		ORDER BY p.created_at DESC`

	results, err := localDatabase.RunDatabaseQuery(ctx, fetchAllCommentsQuery)
	if err != nil {
		return nil, err
	}

	var comments []Comment

	// ce qu'on veut renvoyer
	for _, row := range results {
		comments = append(comments, fromMap(row))
	}

	return comments, nil
}

func FetchComment(ctx context.Context, params map[string]any) ([]Comment, error) {

	comment_ID, comment_IDOK := params["comment_id"].(string)
	post_UUID, post_UUIDOK := params["post_uuid"].(string)
	user_UUID, user_UUIDOK := params["user_uuid"].(string)

	var fetchCommentquery string

	// Param car on souhaite récuperer de la data depuis la db
	var param string

	// Préparer les requetes SQL et fetch si user_uuid ou post_uuid
	if post_UUIDOK {

		fetchCommentquery = `
    SELECT c.comment_id, c.content, c.post_uuid, u.user_uuid, u.username, u.profile_picture, c.created_at
    FROM Comments AS c
    JOIN Users AS u ON c.user_uuid = u.user_uuid
    WHERE c.post_uuid = ?`
		param = post_UUID

	} else if user_UUIDOK {

		fetchCommentquery = `SELECT * FROM Comments WHERE user_uuid = ?`
		param = user_UUID
	} else if comment_IDOK {

		fetchCommentquery = `SELECT * FROM Comments WHERE comment_id = ?`
		param = comment_ID
	} else {
		return nil, errors.New("informations manquantes")
	}

	results, err := localDatabase.RunDatabaseQuery(ctx, fetchCommentquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var comments []Comment

	// Lecture des résultats à partir du slice de maps

	for _, row := range results {
		comments = append(comments, fromMap(row))
	}

	return comments, nil

}

func DeleteComment(ctx context.Context, params map[string]any) error {
	comment_ID, comment_IDOK := params["comment_id"].(int)

	if !comment_IDOK {
		return errors.New("informations manquantes")
	}

	// Exécution de la requête de suppression
	deleteCommentQuery := `DELETE FROM Comments WHERE comment_id = ?`
	_, err := localDatabase.RunDatabaseQuery(ctx, deleteCommentQuery, comment_ID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du commentaire: %v", err)
	}

	return nil
}
