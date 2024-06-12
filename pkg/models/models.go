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
	ID        int       `json:"id,omitempty"`
	Title     string    `json:"title"`
	Params    []*Params `json:"params"`
	Center    *Coord    `json:"center,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	ExpiresIn int       `json:"expiresIn"`
	Expires   time.Time `json:"expires,omitempty"`
}

func NewUniverse(title string, params []*Params, expires int) *Universe {
	return &Universe{
		Title:   title,
		Params:  params,
		Expires: time.Now().AddDate(0, 0, expires),
	}
}

type Params struct {
	Coord *Coord  `json:"coord"`
	Mass  float64 `json:"mass"`
}

type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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
