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

func NewUniverse(title string, params []*Params, expires int) *Universe {
	return &Universe{
		Title:   title,
		Params:  params,
		Expires: time.Now().AddDate(0, 0, expires),
	}
}

type Params struct {
	Coord *Coord
	Mass  float64
}

type Coord struct {
	X, Y, Z float64
}

func NewParams(x, y, z, mass float64) *Params {
	return &Params{
		Coord: &Coord{
			X: x,
			Y: y,
			Z: z,
		},
		Mass: mass,
	}
}
