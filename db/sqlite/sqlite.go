package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func dbConnect(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	// PRAGMAs apply to a connection, so retain one connection for this serial suite.
	db.SetMaxOpenConns(1)
	for _, statement := range []string{
		`PRAGMA synchronous = OFF;`,
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA cache_size = -2097152;`,
		`PRAGMA temp_store = MEMORY;`,
		`PRAGMA locking_mode = EXCLUSIVE;`,
		`PRAGMA mmap_size = 8589934592;`,
	} {
		if _, err := db.Exec(statement); err != nil {
			db.Close()
			return fmt.Errorf("configure SQLite: %w", err)
		}
	}
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
