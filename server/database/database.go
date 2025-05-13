package localDatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

// Variable globale pour la base de données
var Db *sql.DB

// init() est appelé automatiquement avant le main() afin de vérifier la connexion à la db
func InitConnection() {
	var err error

	// On s'assure d'abord qu'une connexion n'existe pas déjà
	if Db != nil {
		if err = Db.Ping(); err == nil {
			return
		}
	}

	// Ouvre la connexion à la base de données
	Db, err = sql.Open("sqlite3", "./realTimeDB.db")
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}

	// Vérifie la connexion
	if err = Db.Ping(); err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données réussie")

	// Crée les tables si elles n'existent pas
	createTables()
}

func createTables() {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS Users (
        user_uuid VARCHAR(36) PRIMARY KEY NOT NULL,
        username TEXT UNIQUE,
        email TEXT UNIQUE,
        password TEXT,
        first_name TEXT,
        last_name TEXT,
        birth_date TIME,
        gender TEXT,
        role TEXT,
        profile_picture MEDIUMTEXT,
        created_at TIME
    );`

	createPostsTable := `
    CREATE TABLE IF NOT EXISTS Posts (
        post_uuid VARCHAR(36) PRIMARY KEY NOT NULL,
        user_uuid VARCHAR(36) NOT NULL,
        title TEXT,
        content TEXT,
        username TEXT,
        profile_picture MEDIUMTEXT,
        likes INT,
        dislikes INT,
        created_at TIME
    );`

	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS Comments (
        comment_id VARCHAR(36) PRIMARY KEY NOT NULL,
        post_uuid  VARCHAR(36) NOT NULL,
        user_uuid  VARCHAR(36) NOT NULL,
        content TEXT,
        username TEXT,
        profile_picture MEDIUMTEXT,
        created_at TIME
    );`

	// Exécute les requêtes pour créer les tables
	queries := []string{createUsersTable, createPostsTable, createCommentsTable}
	for _, query := range queries {
		_, err := Db.Exec(query)
		if err != nil {
			log.Fatalf("Erreur lors de la création des tables : %v", err)
		}
	}

	log.Println("Tables créées ou déjà existantes")
}

// runQuery exécute une requête SQL avec des paramètres et renvoie les résultats

/*
*---------------------------------------------------------------------------------------------------------*
|     																									  |
| - variable query = représente les requetes sql														  |
| - params ...any : le second paramètre est un argument de type variadique (...any),      				  |
| ce qui signifie que l'on peut passer un nombre quelconque d'arguments de différents types (any		  |
| peut être n'importe quel type en Go). Cela permet de fournir dynamiquement les valeurs pour les ?       |
| de la requête SQL (SELECT * FROM users WHERE id = ?).													  |
|																										  |
| []map[string]any : la fonction retourne une slice ([]) de maps (map[string]any),						  |
| où chaque map représente une row de la table avec des columns et leurs valeurs.					      |
|																										  |
*---------------------------------------------------------------------------------------------------------
*/
func RunDatabaseQuery(ctx context.Context, query string, params ...any) ([]map[string]any, error) {
	// Désactiver la validation si les requêtes préparées sont utilisées
	// ou ajuster le regex pour éviter les faux positifs
	re := regexp.MustCompile(`(?i)<[^>]+>|<script.*?>.*?</script.*?>`)

	for _, value := range params {
		valueToRead, ok := value.(string)
		if !ok {
			continue
		}

		if re.FindAllString(valueToRead, -1) != nil {
			return nil, fmt.Errorf("injection detected")
		}
	}

	rows, err := Db.QueryContext(ctx, query, params...)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]any
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Erreur lors du scan des résultats : %v", err)
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			row[col] = val
		}
		results = append(results, row)
	}

	return results, nil
}
