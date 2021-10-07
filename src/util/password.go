package util

import "golang.org/x/crypto/bcrypt"

// HashPassword Create a password hash
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return string(hash), err
	}
	return string(hash), nil
}

// PasswordIsValid Check if the given password is valid
func PasswordIsValid(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}