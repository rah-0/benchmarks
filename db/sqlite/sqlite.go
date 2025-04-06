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
	db.Exec(`PRAGMA synchronous = OFF;`)
	db.Exec(`PRAGMA journal_mode = WAL;`)
	db.Exec(`PRAGMA cache_size = -2097152;`)
	db.Exec(`PRAGMA temp_store = MEMORY;`)
	db.Exec(`PRAGMA locking_mode = EXCLUSIVE;`)
	db.Exec(`PRAGMA mmap_size = 8589934592;`)
	return nil
}

func dbDisconnect() error {
	return db.Close()
}

func dbTableSampleACreate() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS sample_a (
		FieldA TEXT NOT NULL,
		FirstInsert DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(createTableQuery); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_field_a ON sample_a(FieldA);`); err != nil {
		return err
	}

	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_first_insert ON sample_a(FirstInsert);`); err != nil {
		return err
	}

	return nil
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
