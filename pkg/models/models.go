package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Storage interface {
	Insert(title string, x, y, z float64, mass float64, expires string) (int, error)
	Get(id int) (*Object, error)
	Latest() ([]*Object, error)
}

type Object struct {
	ID      int
	Title   string
	X, Y, Z float64
	Mass    float64
	Created time.Time
	Expires time.Time
}
