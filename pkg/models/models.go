package models

import (
	"errors"
	"time"
)

var ErrNotRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: email already used")
var ErrInvalidCredentials = errors.New("models: invalid credentials")

type Snippet struct {
	ID int
	Title string
	Content string
	Expires time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID int
	Name string
	Email string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}