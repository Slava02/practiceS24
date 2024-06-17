package models

import (
	"errors"
	"github.com/Slava02/practiceS24/pkg/validator"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: подходящей записи не найдено")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type UniverseModel interface {
	Insert(obj *Universe) (int, error)
	Get(id int) (*Universe, error)
	Latest(num int) ([]*Universe, error)
}

type Universe struct {
	ID        int       `json:"id,omitempty"`
	Title     string    `json:"title" validate:"required"`
	Params    []*Params `json:"params" validate:"required"`
	Center    *Coord    `json:"center,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	ExpiresIn int       `json:"expiresIn" validate:"required"`
	Expires   time.Time `json:"expires,omitempty"`
	validator.Validator
}

func NewUniverse(title string, params []*Params, expires int) *Universe {
	return &Universe{
		Title:   title,
		Params:  params,
		Expires: time.Now().AddDate(0, 0, expires),
	}
}

type Params struct {
	Coord *Coord  `json:"coord" validate:"required"`
	Mass  float64 `json:"mass" validate:"required"`
}

type Coord struct {
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
	Z float64 `json:"z" validate:"required"`
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
