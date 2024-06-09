package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Storage interface {
	Insert(obj *Universe) (int, error)
	Get(id int) (*Universe, error)
	Latest(num int) ([]*Universe, error)
}

type Universe struct {
	ID      int
	Title   string
	Params  []*Params
	Created time.Time
	Expires time.Time
}

type Params struct {
	Coord *Coord
	Mass  float64
}

type Coord struct {
	X, Y, Z float64
}
