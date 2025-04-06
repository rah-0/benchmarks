package postgresql

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	_ "github.com/lib/pq"

	"github.com/rah-0/benchmarks/util"
)

var (
	db *sql.DB
)

func dbTableSampleAInsert(fieldAValue string) error {
	stmt, err := db.Prepare(`
		INSERT INTO sample_a (FieldA)
		VALUES ($1);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fieldAValue)
	return err
}

func dbConnect() (err error) {
	ip := util.GetEnvVariable("Postgres_IP")
	port := util.GetEnvVariable("Postgres_Port")
	dbName := util.GetEnvVariable("Postgres_Name")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres dbname=%s sslmode=disable", ip, port, dbName)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(runtime.NumCPU())
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return nil
}

func dbDisconnect() error {
	return db.Close()
}

func dbTableSampleACreate() error {
	query := `
	CREATE TABLE IF NOT EXISTS sample_a (
		FieldA VARCHAR(36) NOT NULL DEFAULT '',
		FirstInsert TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
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
