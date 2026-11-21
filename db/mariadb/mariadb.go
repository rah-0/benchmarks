package mariadb

import (
	"database/sql"
	"runtime"
	"time"
)

var (
	db *sql.DB
)

func dbTableSampleAInsert(fieldAValue string) error {
	stmt, err := db.Prepare(`
		INSERT INTO sample_a (FieldA)
		VALUES (?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement
	_, err = stmt.Exec(fieldAValue)
	if err != nil {
		return err
	}

	return nil
}

func dbConnect(dsn string) (err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(runtime.NumCPU())
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
	FieldA VARCHAR(36) NOT NULL DEFAULT '',
	FirstInsert DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	INDEX idx_field_a (FieldA),
	INDEX idx_first_insert (FirstInsert)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPACT
`
	_, err := db.Exec(createTableQuery)
	return err
}

func dbTableSampleADrop() error {
	dropTableQuery := `DROP TABLE IF EXISTS sample_a;`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		return err
	}
	return nil
}
