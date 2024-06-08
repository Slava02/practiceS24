package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/pkg/models"
)

type ObjectModel struct {
	DB *sql.DB
}

func (m *ObjectModel) Insert(obj *models.Object) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: %w", err)
	}

	stmt := `INSERT INTO objects (title, created, expires) 
	VALUES (?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	o, err := m.DB.Exec(stmt, obj.Title, obj.Expires)
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: %w", err)
	}

	id, err := o.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: CAN'T GET LAST INSERT ID: %w", err)
	}

	stmt = `INSERT INTO params (id, x, y, z, mass) 
			VALUES (?, ?, ?, ?, ?)`

	for _, param := range obj.Params {
		_, err = m.DB.Exec(stmt, id, param.Coord.X, param.Coord.Y, param.Coord.Z, param.Mass)
		if err != nil {
			return 0, fmt.Errorf("CAN'T PROCEED INSERT: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT: %w", err)
	}

	return int(id), nil
}

func (m *ObjectModel) Get(id int) (*models.Object, error) {
	o := &models.Object{}

	tx, err := m.DB.Begin()

	stmt := `SELECT id, title, created, expires 
			 FROM objects 
			 WHERE expires > NOW() AND id = ?`

	err = m.DB.QueryRow(stmt, id).Scan(&o.ID, &o.Title, &o.Created, &o.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("CAN'T GET OBJECT: %w", err)
		}
	}

	stmt = `SELECT x, y, z, mass FROM params WHERE id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("CAN'T GET OBJECT: %w", err)
	}

	for rows.Next() {
		var p *models.Params
		if err = rows.Scan(&p.Coord.X, &p.Coord.X, &p.Coord.X, &p.Mass); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, fmt.Errorf("CAN'T GET OBJECT: %w", err)
			}
		}
		o.Params = append(o.Params, p)
	}
	if err = rows.Err(); err != nil {
		return o, err
	}

	err = tx.Commit()

	return o, nil
}

func (m *ObjectModel) Latest(num int) ([]*models.Object, error) {
	stmt1 := `SELECT id, title, created, expires 
			 FROM objects 
			 WHERE expires > NOW() ORDER BY created DESC LIMIT ?`

	stmt2 := `SELECT x, y, z, mass
			 FROM params 
			 WHERE id = ?`

	rows_o, err := m.DB.Query(stmt1, num)
	if err != nil {
		return nil, err
	}
	defer rows_o.Close()

	var objects []*models.Object

	for rows_o.Next() {
		o := &models.Object{}
		err = rows_o.Scan(&o.ID, &o.Title, &o.Created, &o.Expires)
		if err != nil {
			return nil, err
		}

		rows_p, err := m.DB.Query(stmt2)
		if err != nil {
			return nil, err
		}

		for rows_p.Next() {
			var p *models.Params

			err = rows_p.Scan(&p.Coord.X, &p.Coord.Y, &p.Coord.Z, &p.Mass)
			if err != nil {
				return nil, err
			}

			o.Params = append(o.Params, p)
		}

		objects = append(objects, o)
		rows_p.Close()
	}

	if err = rows_o.Err(); err != nil {
		return nil, err
	}

	return objects, nil
}
