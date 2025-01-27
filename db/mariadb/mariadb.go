package mariadb

import (
	"database/sql"
	"runtime"
	"time"

	"github.com/rah-0/benchmarks/util"
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

func dbConnect() (err error) {
	usr := util.GetEnvVariable("MariaDB_Usr")
	pwd := util.GetEnvVariable("MariaDB_Pwd")
	ip := util.GetEnvVariable("MariaDB_IP")
	port := util.GetEnvVariable("MariaDB_Port")
	dbName := util.GetEnvVariable("MariaDB_Name")

	dsn := usr + ":" + pwd + "@tcp(" + ip + ":" + port + ")/" + dbName
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
		Id INT(11) NOT NULL AUTO_INCREMENT,
		FirstInsert DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
		LastUpdate DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
		FieldA VARCHAR(36) NOT NULL DEFAULT '',
		PRIMARY KEY (Id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPACT`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}

func dbTableSampleADrop() error {
	dropTableQuery := `DROP TABLE IF EXISTS sample_a;`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		return err
	}
	return nil
}
