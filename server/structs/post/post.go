package post

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	localDatabase "github.com/RealTimeForumDosJean/server/database"
	"github.com/RealTimeForumDosJean/server/structs/user"
)

type Post struct {
	UUID            string    `json:"post_uuid"`
	User_uuid       string    `json:"user_uuid"`
	Username        string    `json:"username"`
	Profile_picture string    `json:"profile_picture"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Likes           int       `json:"likes"`
	Dislikes        int       `json:"dislikes"`
	Created_at      time.Time `json:"created_at"`
}

func FromMap(values map[string]any) Post {
	newPost := Post{}

	if post_uuid, uuidOK := values["post_uuid"].(string); uuidOK {
		newPost.UUID = post_uuid
	}
	if user_uuid, uuidOK := values["user_uuid"].(string); uuidOK {
		newPost.User_uuid = user_uuid
	}
	if username, usernameOK := values["username"].(string); usernameOK {
		newPost.Username = username
	}
	if profilePic, pictureOK := values["profile_picture"].(string); pictureOK {
		newPost.Profile_picture = profilePic
	}
	if title, uuidOK := values["title"].(string); uuidOK {
		newPost.Title = title
	}
	if content, uuidOK := values["content"].(string); uuidOK {
		newPost.Content = content
	}
	if likes, uuidOK := values["likes"].(int64); uuidOK {
		newPost.Likes = int(likes)
	}
	if dislikes, uuidOK := values["dislikes"].(int64); uuidOK {
		newPost.Dislikes = int(dislikes)
	}
	if created_at, uuidOK := values["created_at"].(time.Time); uuidOK {
		newPost.Created_at = created_at
	}

	return newPost
}

func (p *Post) ToMap() map[string]any {
	postMap := make(map[string]any, 0)

	postMap["post_uuid"] = p.UUID
	postMap["user_uuid"] = p.User_uuid
	postMap["username"] = p.Username
	postMap["profile_picture"] = p.Profile_picture
	postMap["title"] = p.Title
	postMap["content"] = p.Content
	postMap["likes"] = p.Likes
	postMap["dislikes"] = p.Dislikes
	postMap["created_at"] = p.Created_at.Format("2006-01-02T15:04:05Z")

	return postMap
}

func CreatePost(r *http.Request, params map[string]any) error {

	results, err := localDatabase.RunDatabaseQuery(
		r.Context(),
		"SELECT username, profile_picture FROM Users WHERE user_uuid = ?",
		params["user_uuid"])
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des infos utilisateur: %v", err)
	}

	usrFound := user.FromMap(results[0])

	createPostQuery := `INSERT INTO Posts (post_uuid, user_uuid, username, profile_picture, title, content, likes, dislikes, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = localDatabase.RunDatabaseQuery(r.Context(), createPostQuery, params["post_uuid"], params["user_uuid"], usrFound.Username, usrFound.ProfilePicture, params["title"], params["content"], int64(0), int64(0), params["created_at"])
	if err != nil {
		return fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	return nil
}

func FetchPost(ctx context.Context, params map[string]any) ([]Post, error) {
	var fetchPostquery string
	var param string

	fetchPostquery = `
		SELECT p.*, u.username, u.profile_picture
		FROM Posts p
		JOIN Users u ON p.user_uuid = u.user_uuid
		WHERE p.`

	if post_UUID, ok := params["post_uuid"].(string); ok {
		fetchPostquery += "post"
		param = post_UUID
	} else if user_UUID, ok := params["user_uuid"].(string); ok {
		fetchPostquery += "user"
		param = user_UUID
	} else {
		return nil, errors.New("informations manquantes")
	}

	fetchPostquery += "_uuid = ?"

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchPostquery, param)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	var posts []Post

	for _, row := range rows {
		posts = append(posts, FromMap(row))
	}

	return posts, nil
}

func FetchAllPosts(ctx context.Context) ([]Post, error) {
	fetchAllPostsQuery := `
        SELECT p.*, u.username, u.profile_picture 
        FROM Posts p
        JOIN Users u ON p.user_uuid = u.user_uuid
        ORDER BY p.created_at DESC`

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchAllPostsQuery)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var posts []Post
	for _, row := range rows {
		posts = append(posts, FromMap(row))
	}

	return posts, nil
}

func DeletePost(ctx context.Context, params map[string]any) error {
	post_UUID, post_UUIDOK := params["post_uuid"].(string)

	if !post_UUIDOK {
		return errors.New("informations manquantes")
	}

	deletePostQuery := `DELETE FROM Posts WHERE post_uuid = ?`
	_, err := localDatabase.RunDatabaseQuery(ctx, deletePostQuery, post_UUID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du post: %v", err)
	}

	return nil
}

func FetchAllCategories(ctx context.Context) ([]string, error) {
	fetchAllCategoriesQuery := `
        SELECT DISTINCT categories
        FROM Posts
        WHERE categories IS NOT NULL AND categories <> ''`

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchAllCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	var categories []string
	for _, row := range rows {
		// On récupère la colonne 'categories' et on la divise par les virgules
		if catStr, ok := row["categories"].(string); ok {
			for _, category := range strings.Split(catStr, ",") {
				// Éliminer les espaces superflus et vérifier si la catégorie n'est pas déjà présente
				trimmedCategory := strings.TrimSpace(category)
				if trimmedCategory != "" {
					categories = append(categories, "#"+trimmedCategory)
				}
			}
		}
	}

	// Utiliser un map pour éliminer les doublons
	categoryMap := make(map[string]struct{})
	for _, category := range categories {
		categoryMap[category] = struct{}{}
	}

	// Convertir le map en slice
	uniqueCategories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		uniqueCategories = append(uniqueCategories, category)
	}

	return uniqueCategories, nil
}

// Mini algo pour trouver les tendances

/*----------------------------------------------------------------------------------------------------

Explication de la Fonction
Requête SQL :

La requête utilise une sous-requête pour séparer les catégories (hashtags) en lignes distinctes.
SUBSTRING_INDEX et numbers génèrent les catégories individuelles à partir d'une chaîne qui contient plusieurs catégories séparées par des virgules.
Le nombre d'occurrences de chaque catégorie est ensuite compté et celles ayant plus de 10 occurrences sont sélectionnées.
Traitement des Résultats :

La fonction retourne un map[string]int où la clé est le nom de la catégorie (hashtag) et la valeur est le nombre de posts qui utilisent cette catégorie.

----------------------------------------------------------------------------------------------------*/

func FetchCategoryRanking(ctx context.Context) (map[string]int, error) {
	// Requête pour compter le nombre d'occurrences de chaque catégorie
	fetchCategoryRankingQuery := `
        SELECT category, COUNT(*) AS count
        FROM (
            SELECT TRIM(SUBSTRING_INDEX(SUBSTRING_INDEX(categories, ',', numbers.n), ',', -1)) AS category
            FROM Posts
            INNER JOIN (
                SELECT 1 AS n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL 
                SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9 UNION ALL SELECT 10
            ) AS numbers ON CHAR_LENGTH(categories) - CHAR_LENGTH(REPLACE(categories, ',', '')) >= numbers.n - 1
        ) AS subquery
        WHERE category <> ''
        GROUP BY category
        HAVING count > 10
        ORDER BY count DESC`

	// Exécute la requête
	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchCategoryRankingQuery)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du classement des catégories: %v", err)
	}

	// Préparez le mappage des catégories et leur compte
	categoryRanking := make(map[string]int)
	for _, row := range rows {
		if category, ok := row["category"].(string); ok {
			if count, ok := row["count"].(int64); ok {
				categoryRanking[category] = int(count)
			}
		}
	}

	return categoryRanking, nil
}
