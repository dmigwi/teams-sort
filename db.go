package main

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createdb = `
		CREATE TABLE IF NOT EXISTS football (
			id UUID NOT NULL,
			div VARCHAR NOT NULL,
			date_time TIMESTAMP NOT NULL,
			home_team VARCHAR NOT NULL,
			away_team VARCHAR NOT NULL,
			fthg INT DEFAULT(0),
			ftag INT DEFAULT(0),
			CONSTRAINT football_pk PRIMARY KEY (id)
		);`

	insertQry = `INSERT INTO football VALUES (?,?,?,?,?,?,?);`
	readQry   = `SELECT * FROM football LIMIT %d;`
	countQry  = `SELECT COUNT(*) FROM football;`
)

type dbConn struct {
	db *sql.DB
}

// setUpDatabase sets the correct db location path depending on the path provided.
// Initialises the database connection and sets up the football table if it doesn't
// exist.
func setUpDatabase(dbLocation string) (*dbConn, error) {
	path := "." // default to current location if path ot provided
	if dbLocation != "" {
		path = filepath.Dir(dbLocation)
	}

	dbLocation = fmt.Sprintf("%s/teams_sort.db", path)
	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return nil, errors.New("failed to initialise a db connection :" + err.Error())
	}

	// Set up football table if it doesn't exist.
	if _, err := db.Exec(createdb); err != nil {
		return nil, errors.New("failed to create table football :" + err.Error())
	}

	return &dbConn{db}, nil
}

// insertData populates the provided data into that database.
func (conn *dbConn) insertData(records []matchInfo) error {
	if conn.db == nil {
		return errors.New("missing database connection")
	}

	isExist, err := conn.isDbPrepopulated()
	if err != nil {
		return err
	}

	// Avoid duplicating the same data everytime data is loaded.
	if !isExist {
		for _, r := range records {
			uuidVal, err := uuid.NewUUID()
			if err != nil {
				return errors.New("failed to create uuid :" + err.Error())
			}
			_, err = conn.db.Exec(insertQry, uuidVal.String(), r.Division, r.Date, r.HomeTeam, r.AwayTeam, r.FTHG, r.FTAG)
			if err != nil {
				return errors.New("failed to insert data :" + err.Error())
			}
		}
	}

	return nil
}

// readData returns the number of records with the provided limit starting from
// the most recent.
func (conn *dbConn) readData(limit int) ([]matchInfo, error) {
	if conn.db == nil {
		return nil, errors.New("missing database connection")
	}

	rows, err := conn.db.Query(fmt.Sprintf(readQry, limit))
	if err != nil {
		return nil, errors.New("failed to read data :" + err.Error())
	}

	data := make([]matchInfo, 0)
	for rows.Next() {
		r := matchInfo{}
		if err = rows.Scan(&r.ID, &r.Division, &r.Date, &r.HomeTeam, &r.AwayTeam, &r.FTHG, &r.FTAG); err != nil {
			return nil, errors.New("failed to scan row data :" + err.Error())
		}
		data = append(data, r)
	}
	return data, nil
}

// isDbPrepopulated returns true if the database has been prepopulated.
func (conn *dbConn) isDbPrepopulated() (bool, error) {
	var rowsCount int
	row := conn.db.QueryRow(countQry)
	if err := row.Scan(&rowsCount); err != nil {
		return false, errors.New("failed to get the rows count :" + err.Error())
	}
	return rowsCount > 0, nil
}
