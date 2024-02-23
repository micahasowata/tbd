package utils

import "testing"

func TestHashPassword(t *testing.T) {
	plaintext := "password"
	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// Length of a normal bcrypt hash is 60
	if len(hash) != 60 {
		t.Error("hash should be 60 chracters long")
	}

	if string(hash) == plaintext {
		t.Error("hash should not be the same as plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	plaintext := "password"
	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	err = CheckPassword(hash, plaintext)
	if err != nil {
		t.Errorf("password verification failed %s", err.Error())
	}
}
