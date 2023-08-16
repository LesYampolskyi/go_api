package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 10
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
	// maxFirstNameLen = 12
	// maxLastNameLen = 12
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func (params CreateUserParams) Validate() []string {
	errors := []string{}
	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("First name length should be at least %d characters.", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("Last name length should be at least %d characters.", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("Password name length should be at least %d characters.", minPasswordLen))
	}
	if !isEmailValid(params.Email) {
		errors = append(errors, fmt.Sprintf("Email is invalid"))
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
