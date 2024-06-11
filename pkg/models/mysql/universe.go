package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Slava02/practiceS24/pkg/models"
	"math"
)

type UniverseModel struct {
	DB *sql.DB
}

func (m *UniverseModel) Insert(obj *models.Universe) (int, error) {

	//interval := strconv.Itoa(RoundTime(obj.Expires.Sub(time.Now()).Seconds() / 86400))

	stmt := `INSERT INTO objects (title, created, expires) 
	VALUES (?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	o, err := m.DB.Exec(stmt, obj.Title, obj.ExpiresIn)
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT (1): %w", err)
	}

	id, err := o.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CAN'T PROCEED INSERT (2): CAN'T GET LAST INSERT ID: %w", err)
	}

	stmt = `INSERT INTO params (id, x, y, z, mass) 
			VALUES (?, ?, ?, ?, ?)`

	for _, param := range obj.Params {
		_, err = m.DB.Exec(stmt, id, param.Coord.X, param.Coord.Y, param.Coord.Z, param.Mass)
		if err != nil {
			return 0, fmt.Errorf("CAN'T PROCEED INSERT (3): %w", err)
		}
	}

	return int(id), nil
}

func (m *UniverseModel) Get(id int) (*models.Universe, error) {
	o := &models.Universe{}

	//tx, err := m.DB.Begin()

	stmt := `SELECT id, title, created, expires 
			 FROM objects 
			 WHERE expires > NOW() AND id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&o.ID, &o.Title, &o.Created, &o.Expires)
	// log.Printf("INFO: Got object fields: %+v", o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, fmt.Errorf("CAN'T GET OBJECT (1): %w", err)
		}
	}

	stmt = `SELECT x, y, z, mass FROM params WHERE id = ?`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, fmt.Errorf("CAN'T GET OBJECT (2): %w", err)
	}

	for rows.Next() {
		p := &models.Params{
			Coord: &models.Coord{},
		}
		err = rows.Scan(&p.Coord.X, &p.Coord.X, &p.Coord.X, &p.Mass)
		// log.Printf("INFO: Got params: %+v\n", p)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrNoRecord
			} else {
				return nil, fmt.Errorf("CAN'T GET OBJECT (3): %w", err)
			}
		}
		o.Params = append(o.Params, p)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("CAN'T GET OBJECT (4): %w", err)
	}

	//err = tx.Commit()
	// log.Printf("INFO: Send object: %+v\n", o)
	return o, nil
}

func (m *UniverseModel) Latest(num int) ([]*models.Universe, error) {
	stmt1 := `SELECT id, title, created, expires 
			 FROM objects 
			 WHERE expires > NOW() ORDER BY created DESC LIMIT ?`

	stmt2 := `SELECT x, y, z, mass
			 FROM params 
			 WHERE id = ?`

	rows_o, err := m.DB.Query(stmt1, num)
	if err != nil {
		return nil, fmt.Errorf("CAN'T GET LATEST (1): %w", err)
	}
	defer rows_o.Close()

	var objects []*models.Universe

	for rows_o.Next() {
		o := &models.Universe{}
		err = rows_o.Scan(&o.ID, &o.Title, &o.Created, &o.Expires)
		// log.Printf("INFO: Got object values: %+v\n", o)
		if err != nil {
			return nil, fmt.Errorf("CAN'T GET LATEST (2): %w", err)
		}

		rows_p, err := m.DB.Query(stmt2, o.ID)
		if err != nil {
			return nil, fmt.Errorf("CAN'T GET LATEST (3): %w", err)
		}

		for rows_p.Next() {
			p := &models.Params{
				Coord: &models.Coord{},
			}

			err = rows_p.Scan(&p.Coord.X, &p.Coord.Y, &p.Coord.Z, &p.Mass)
			// log.Printf("INFO: Got params: %+v\n", p)
			if err != nil {
				return nil, fmt.Errorf("CAN'T GET LATEST (4): %w", err)
			}

			o.Params = append(o.Params, p)
		}

		objects = append(objects, o)
		rows_p.Close()
	}

	if err = rows_o.Err(); err != nil {
		return nil, fmt.Errorf("CAN'T GET LATEST (5): %w", err)
	}

	return objects, nil
}

func RoundTime(input float64) int {
	var result float64
	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}
	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)
	return int(i)
}
