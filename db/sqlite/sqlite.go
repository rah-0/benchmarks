package sqlite

import (
	"database/sql"
	"runtime"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func dbConnect() error {
	var err error
	db, err = sql.Open("sqlite3", "./benchmark.db")
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(runtime.NumCPU())
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 1)
	return nil
}

func dbDisconnect() error {
	return db.Close()
}

func dbTableSampleACreate() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS sample_a (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		FirstInsert DATETIME DEFAULT CURRENT_TIMESTAMP,
		LastUpdate DATETIME DEFAULT CURRENT_TIMESTAMP,
		FieldA TEXT NOT NULL
	);
	`
	_, err := db.Exec(createTableQuery)
	return err
}

func dbTableSampleADrop() error {
	_, err := db.Exec(`DROP TABLE IF EXISTS sample_a;`)
	return err
}

func dbTableSampleAInsert(fieldAValue string) error {
	stmt, err := db.Prepare(`INSERT INTO sample_a (FieldA) VALUES (?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fieldAValue)
	return err
}
