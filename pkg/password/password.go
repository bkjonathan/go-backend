package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("hashing password with bcrypt: %w", err)
	}

	return string(hash), nil
}

func Compare(hash, raw string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)); err != nil {
		return fmt.Errorf("comparing password hash: %w", err)
	}
	return nil
}
