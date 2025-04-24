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

	// Ouvre la connexion à la base de données ici
	// Le fichier realTimeDB.db est à la racine, c'est ce fichier qui contient toute la database
	// Grâce à l'extension sqlite de vscode, nous pouvons visualiser cela plus facilement
	// Clic droit sur realTimeDB.db
	// Open database
	// Magie on peut voir les tables avec les columns et rows

	Db, err = sql.Open("sqlite3", "./realTimeDB.db")
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données : %v", err)
	}

	// Vérifie la connexion
	if err = Db.Ping(); err != nil {
		log.Fatalf("Erreur lors de la connexion à la base de données : %v", err)
	}

	log.Println("Connexion à la base de données réussie")
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

	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	for _, value := range params {
		valueToRead, ok := value.(string)
		if !ok {
			continue
		}

		if re.FindAllString(valueToRead, -1) != nil {
			return nil, fmt.Errorf("injection detected")
		}
	}

	//----------------------------------------------------------------------//
	// Prépare la requête
	// Voir ./database.go "var Db *sql.DB"
	// params ex: "SELECT * FROM users"
	//----------------------------------------------------------------------//

	rows, err := Db.QueryContext(ctx, query, params...)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	defer rows.Close()

	// Récupère les colonnes de la requête
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Création d'une slice pour stocker les valeurs

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Stockage des résultats dans une liste de maps
	// Tous les éléments trouvés sont stockés et renvoyés

	var results []map[string]any
	for rows.Next() {
		// Remplit les valeurs pour la ligne actuelle
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Erreur lors du scan des résultats : %v", err)
			return nil, err
		}

		// Crée une map pour la ligne
		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			row[col] = val
		}
		results = append(results, row)
	}

	return results, nil
}
