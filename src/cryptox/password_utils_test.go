package cryptox

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "12i819jcasd13"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
		return
	}
	ok := CheckPasswordHash(password, hash)
	if !ok {
		t.Error("wrong password")
		return
	}

	newPassword := "dreissig"
	newHash := "$2a$14$bp52QfoU6kf6n3QOZzKeheCuL71FrCc.GDFm7mVVIpmLAs1D6WwVy"
	ok = CheckPasswordHash(newPassword, newHash)
	if !ok {
		t.Error("wrong password of given")
		return
	}
}
