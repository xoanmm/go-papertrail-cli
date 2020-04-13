package papertrail

import (
	"errors"
	"fmt"
	"time"
)

const shortDateFormat = "01/02/2006 15:04:05"

func GetTimeStampUnixFromDate(date string) (int64, error) {
	layout := "01/02/2006 15:04:05"
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0, errors.New("Error converting date " + date + "to timestamp ")
	}
	return t.Unix(), nil
}

func GetTimeInUTCFromUnixTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(shortDateFormat)
}

//ShortDateFromString parse shot date from string
func ShortDateFromString(ds string) (time.Time, error) {
	t, err := time.Parse(shortDateFormat, ds)
	if err != nil {
		return t, err
	}
	return t, nil
}

//CheckDataBoundariesStr checks is startdate <= enddate
func CheckDataBoundariesStr(startdate, enddate string) (bool, error) {

	tstart, err := ShortDateFromString(startdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := ShortDateFromString(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("startdate > enddate - please set proper data boundaries")
	}
	return true, err
}
