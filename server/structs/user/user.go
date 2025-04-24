package user

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	localDatabase "github.com/RealTimeForumDosJean/server/database"
	"golang.org/x/crypto/bcrypt"
)

const (
	SEPARATOR = "|"
)

type User struct {
	UUID              string    `json:"user_uuid"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	EncryptedPassword string    `json:"password"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	BirthDate         time.Time `json:"birth_date"`
	Gender            string    `json:"gender"`
	Contacts          string    `json:"contacts"`

	Role           string    `json:"role"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func FromMap(row map[string]any) User {
	newUser := User{}

	// Utiliser des assertions de type avec vérification de valeur nulle
	if v, ok := row["user_uuid"]; ok && v != nil {
		newUser.UUID = v.(string)
	}
	if v, ok := row["username"]; ok && v != nil {
		newUser.Username = v.(string)
	}
	if v, ok := row["email"]; ok && v != nil {
		newUser.Email = v.(string)
	}
	if v, ok := row["password"]; ok && v != nil {
		newUser.EncryptedPassword = v.(string) // On a vraiment besoin de récupérer le MDP ?
	}
	if v, ok := row["first_name"]; ok && v != nil {
		newUser.FirstName = v.(string)
	}
	if v, ok := row["last_name"]; ok && v != nil {
		newUser.LastName = v.(string)
	}
	if v, ok := row["birth_date"]; ok && v != nil {
		parsedTime, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, "dommage")
		}
		newUser.BirthDate = parsedTime
	}
	if v, ok := row["gender"]; ok && v != nil {
		newUser.Gender = v.(string)
	}
	if v, ok := row["contacts"]; ok && v != nil {
		newUser.Contacts = v.(string)
	}
	if v, ok := row["role"]; ok && v != nil {
		newUser.Role = v.(string)
	}
	if v, ok := row["profile_picture"]; ok && v != nil {
		newUser.ProfilePicture = v.(string)
	}
	if v, ok := row["created_at"]; ok && v != nil {
		parsedTime, err := time.Parse("2006-01-02", v.(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, "dommage")
		}
		newUser.CreatedAt = parsedTime
	}

	return newUser
}

func (u *User) ToMap() map[string]any {
	usrMap := make(map[string]any, 0)

	usrMap["user_uuid"] = u.UUID
	usrMap["username"] = u.Username
	usrMap["email"] = u.Email
	usrMap["password"] = u.EncryptedPassword
	usrMap["first_name"] = u.FirstName
	usrMap["last_name"] = u.LastName
	usrMap["birth_date"] = u.BirthDate.Format("2006-01-02")
	usrMap["gender"] = u.Gender
	usrMap["role"] = u.Role
	usrMap["profile_picture"] = u.ProfilePicture
	usrMap["created_at"] = u.CreatedAt.Format("2006-01-02")

	return usrMap
}

func (u *User) ToCookieValue() string {
	return u.UUID + SEPARATOR +
		u.Username + SEPARATOR +
		u.Email + SEPARATOR +
		u.Role
}

// Trouver un utilisateur par son email et le renvoyer ( pour login )
func FetchUserByEmail(ctx context.Context, email string) (User, error) {
	fetchUserQuery := `SELECT * FROM Users WHERE email = ?`
	params := []any{email}

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchUserQuery, params...)
	if err != nil {
		return User{}, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même email. C'est normal ça ?")
	} else if len(rows) == 0 {
		return User{}, nil
	}

	result := rows[0]

	return FromMap(result), nil
}

// Trouver un utilisateur par son nom et le renvoyer (pour les DMs)
func FetchUserByName(ctx context.Context, name string) (User, error) {
	fetchUserQuery := `SELECT * FROM Users WHERE username = ?`
	params := []any{name}

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchUserQuery, params...)
	if err != nil {
		return User{}, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même nom. C'est vraiment pas normal.")
	} else if len(rows) == 0 {
		return User{}, nil
	}

	result := rows[0]

	return FromMap(result), nil
}

// Savoir si un utilisateur existe par son nom d'utilisateur ( pour register )
func IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	fetchUserQuery := `SELECT * FROM Users WHERE username= ?`
	params := []any{username}

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchUserQuery, params...)
	if err != nil {
		return false, fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	return len(rows) >= 1, nil
}

// Trouver l'image de profil utilisateur avec ID
func FetchPPByID(ctx context.Context, id string) (string, error) {
	fetchUserQuery := `SELECT profile_picture FROM Users WHERE user_uuid= ?`
	params := []any{id}

	rows, err := localDatabase.RunDatabaseQuery(ctx, fetchUserQuery, params...)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du formulaire: %v", err)
	}

	if len(rows) > 1 {
		fmt.Fprintln(os.Stderr, "Y'a plus d'un user avec le même ID. C'est normal ça ?")
	} else if len(rows) == 0 {
		return "", nil
	}

	usrFound := User{}
	result := rows[0]

	// Utiliser une assertion de type avec vérification de valeur nulle
	if v, ok := result["profile_picture"]; ok && v != nil {
		usrFound.ProfilePicture = v.(string)
	}

	return usrFound.ProfilePicture, nil
}

// Enregistrer un user complet ( Register )
func RegisterUser(ctx context.Context, params map[string]any) error {
	profile_picture, ok := params["profile_picture"].(string)
	if !ok {
		profile_picture = "/static/img/icon-7797704_640.png"
	}

	registerUserQuery := `INSERT INTO Users (user_uuid, username, email, password, first_name, last_name, birth_date, gender, role, profile_picture, created_at)  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var err error

	_, err = localDatabase.RunDatabaseQuery(ctx, registerUserQuery, params["user_uuid"], params["username"], params["email"], params["password"], params["first_name"], params["last_name"], params["birth_date"], params["gender"], params["role"], profile_picture, params["created_at"])

	if err != nil {
		return err
	}

	return nil
}

// Mettre à jour les valeurs d'un utilisateur ( Update )
func (u *User) UpdateUser(ctx context.Context, params map[string]any) error {
	re := regexp.MustCompile(`(?i)<[^>]+>|(SELECT|UPDATE|DELETE|INSERT|DROP|FROM|COUNT|AS|WHERE|--)|^\s|^\s*$|<script.*?>.*?</script.*?>`)

	if params["password"] != "" {
		for key, value := range params {
			if (key == "username" || key == "email" || key == "password") && re.FindAllString(value.(string), -1) != nil {
				return fmt.Errorf("injection detected")
			}
		}
	}

	updateUserQuery := `UPDATE Users SET username = ?, email = ?, password = ?, first_name = ?, last_name = ?, birth_date = ?, gender = ?, role = ?, profile_picture = ? WHERE user_uuid = ?`
	_, err := localDatabase.RunDatabaseQuery(ctx, updateUserQuery, params["username"], params["email"], params["password"], params["first_name"], params["last_name"], params["birth_date"], params["gender"], params["role"], params["profile_picture"], params["user_uuid"])

	if err != nil {
		return err
	}

	return nil
}

// FetchAllComments récupère tous les commentaires de la base de données
func FetchAllUsers(ctx context.Context) ([]User, error) {
	results, err := localDatabase.RunDatabaseQuery(ctx, "SELECT user_uuid, username, profile_picture, created_at FROM Users")
	if err != nil {
		return nil, err
	}

	var users []User

	// ce qu'on veut renvoyer
	for _, row := range results {
		users = append(users, FromMap(row))
	}

	return users, nil
}
