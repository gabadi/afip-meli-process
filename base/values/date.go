package values

import (
	"time"
)

const dateFormat = "2006-01-02T15:04:05.000-03:00"

type Date struct {
	time.Time
}

func (date *Date) UnmarshalCSV(csv string) (err error) {
	parsedTime, err := time.Parse(dateFormat, csv)
	if err != nil {
		return err
	}

	date.Time = parsedTime.In(buenosAiresLocation)
	return nil
}

func getBuenosAiresLocation() *time.Location {
	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		panic(err)
	}
	return loc
}

var buenosAiresLocation = getBuenosAiresLocation()
