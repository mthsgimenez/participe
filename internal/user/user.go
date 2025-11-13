package user

import (
	"strings"

	"github.com/mthsgimenez/participe/internal/company"
	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	ROLE_USER  = iota
	ROLE_ADMIN = iota
)

func (r UserRole) String() string {
	return [...]string{"ROLE_USER", "ROLE_ADMIN"}[r]
}

func StringToUserRole(s string) UserRole {
	switch strings.ToUpper(s) {
	case "ROLE_ADMIN":
		return ROLE_ADMIN
	default:
		return ROLE_USER
	}
}

type User struct {
	Id      int             `json:"id"`
	Email   string          `json:"email"`
	Company company.Company `json:"company"`
	Name    string          `json:"name"`
	Role    UserRole        `json:"role"`
	hash    string          `json:"-"`
}

func (u *User) SetPassword(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	u.hash = string(hash)
	return nil
}

func (u *User) CheckPassword(plaintext string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.hash), []byte(plaintext)); err != nil {
		return false
	}

	return true
}
