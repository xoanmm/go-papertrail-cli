package papertrail

import (
	"errors"
	"time"
)

const shortDateFormat = "01/02/2006 15:04:05"

// GetTimeStampUnixFromDate returns the date provided as a
// parameter in timestamp format to date using the layout
func GetTimeStampUnixFromDate(date string) (int64, error) {
	layout := "01/02/2006 15:04:05"
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0, errors.New("Error converting date " + date + "to timestamp ")
	}
	return t.Unix(), nil
}

// GetTimeInUTCFromUnixTime returns the date provided as unix timestamp
func GetTimeInUTCFromUnixTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(shortDateFormat)
}

// CheckDateFormat check if date complish format specified
func CheckDateFormat(ds string) error {
	_, err := time.Parse(shortDateFormat, ds)
	if err != nil {
		return err
	}
	return nil
}
