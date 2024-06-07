package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/pkg/models"
	"log"
)

type ObjectModel struct {
	DB *sql.DB
}

func (m *ObjectModel) Insert(title string, x, y, z float64, mass float64, expires string) (int, error) {

	stmt := `INSERT INTO objects (title, x_coord, y_coord, z_coord, mass, created, expires) 
	VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, x, y, z, mass, expires)
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: CAN'T GET LAST INSERT ID: %w", err)
	}

	return int(id), nil
}

func (m *ObjectModel) Get(id int) (*models.Object, error) {
	o := &models.Object{}

	stmt := `SELECT id, title, x_coord, y_coord, z_coord, created, expires 
			 FROM objects 
			 WHERE expires > NOW() AND id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&o.ID, &o.Title, &o.X, &o.Y, &o.Z, &o.Created, &o.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("CAN'T GET OBJECT: %w", err)
		}
	}

	log.Printf("GOY OBJECT: title: %s, x: %f,  y: %f,  z: %f,  m: %f", o.Title, o.X, o.Y, o.Z, o.Mass)
	return o, nil
}

func (m *ObjectModel) Latest() ([]*models.Object, error) {
	return nil, nil
}
