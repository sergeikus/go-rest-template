package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"net/http"

	"golang.org/x/crypto/pbkdf2"
)

// Auth represents an interface which wraps different
// authentication scenarios:
//     * Client-Side Session Management (JWT - JSON Web Token)
//     * Server-Side Session Management (Session ID in cookie)
type Auth interface {
	CreateSession(w http.ResponseWriter) (string, error)
	CheckSession(w http.ResponseWriter, r *http.Request) error
	Logout(r *http.Request) error
	PBKDF2HashPassword(password string, salt string) string
}

func pbkdf2HashPassword(password string, salt string, iterations int, keyLenght int) string {
	b := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLenght, sha512.New)
	return fmt.Sprintf("%x", b)
}

const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowercase = "abcdefghijklmnopqrstuvwxyz"
const letters = uppercase + lowercase
const numbers = "1234567890"
const mix = letters + numbers

// GenerateRandomString generates random string
func GenerateRandomString(lenght int) (string, error) {
	bs, err := randomBytes(lenght)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	for i, b := range bs {
		bs[i] = mix[b%byte(len(mix))]
	}
	return string(bs), nil
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
