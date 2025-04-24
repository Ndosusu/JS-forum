package forumServer

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GenerateUUID génère un UUID aléatoire
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = b[6]&^0xf0 | 0x40 // Version 4
	b[8] = b[8]&^0x3f | 0x80 // Variante RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func GetUserUUIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("UserLogged")
	if err != nil {
		return "", errors.New("user not logged in")
	}

	// Décoder la valeur du cookie
	// Supposons que la valeur du cookie soit au format "uuid|username|email"
	parts := strings.Split(cookie.Value, "|")
	if len(parts) < 3 {
		return "", errors.New("invalid cookie format")
	}

	return parts[0], nil // Retourne le user_uuid
}

func IndexOf[T comparable](arr []T, match T) int {
	if len(arr) <= 0 {
		return -1
	}

	for i, elem := range arr {
		if elem == match {
			return i
		}
	}

	return -1
}
