package hyperion

import (
	"time"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/hconn"
)

var (
	connection *hconn.HConn
)

func dbSampleInsert(fieldValue string) error {
	s := SampleV1.Sample{
		Name: fieldValue,
	}
	return s.DbInsert(connection)
}

func dbSampleInsertWithDate(fieldValue string, t time.Time) error {
	s := SampleV1.Sample{
		Name:  fieldValue,
		Birth: t,
	}
	return s.DbInsert(connection)
}
