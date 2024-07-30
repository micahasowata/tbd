package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(plaintext string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
}

func CheckPassword(hash []byte, plaintext string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
}
