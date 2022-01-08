package models

import (
	"errors"
	"time"
)

var ErrNotRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID int
	Title string
	Content string
	Expires time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}