// Package user contains users related CRUD functionality
package user

import (
	"github.com/lib/pq"
	"time"
)

// User represents an individual user
type User struct {
	ID           string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password_hash" json:"password_hash"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new user
type NewUser struct {
	Name            string         `json:"name" validate:"required"`
	Email           string         `json:"email" validate:"required,email"`
	Roles           pq.StringArray `json:"roles" validate:"required"`
	Password        string         `json:"password" validate:"required"`
	PasswordConfirm string         `json:"password_confirm" validate:"required"`
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields, so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// marshalling/unmarshalling
type UpdateUser struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Roles           []string `json:"roles"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}