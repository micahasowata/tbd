package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "a"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	// Length of a normal bcrypt hash is 60
	if len(hashed) != 60 {
		t.Error("Hash should be 60 chracters long")
	}

	if string(hashed) == password {
		t.Error("Hash should not be the same as password")
	}
}
