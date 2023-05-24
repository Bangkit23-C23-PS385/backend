package helper

import (
	"fmt"
	"strconv"
	"strings"
	"ta/backend/src/constant"
	"time"

	"github.com/pkg/errors"
)

func ParseStringToTime(insertedTime string) (parsedTime time.Time, err error) {
	splitted := strings.Split(insertedTime, " ")
	if len(splitted) != 2 {
		err = constant.ErrInvalidFormat
		return
	}

	// Date
	splittedDate := strings.Split(splitted[0], "/")
	if len(splittedDate) != 3 {
		err = constant.ErrInvalidFormat
		return
	}
	day, err := strconv.Atoi(splittedDate[0])
	m, err := strconv.Atoi(splittedDate[1])
	year, err := strconv.Atoi(splittedDate[2])
	if err != nil {
		err = constant.ErrInvalidFormat
		return
	}
	month := time.Month(m)

	// Time
	splittedTime := strings.Split(splitted[1], ":")
	if len(splittedTime) != 2 {
		err = constant.ErrInvalidFormat
		return
	}
	hour, err := strconv.Atoi(splittedTime[0])
	minute, err := strconv.Atoi(splittedTime[1])
	if err != nil {
		err = constant.ErrInvalidFormat
		return
	}

	// Timezone
	location, err := time.LoadLocation(constant.TimeLocation)
	if err != nil {
		err = errors.Wrap(err, "time: load location")
		return
	}

	parsedTime = time.Date(year, month, day, hour, minute, 0, 0, location)

	return
}

func ParseTimeToString(t time.Time) (str string) {
	// Timezone
	location, err := time.LoadLocation(constant.TimeLocation)
	if err != nil {
		err = errors.Wrap(err, "time: load location")
		return
	}

	timeStr := t.In(location).String()
	splittedTimeStr := strings.Split(timeStr, " ")

	date := splittedTimeStr[0]
	splittedDate := strings.Split(date, "-")
	newDate := fmt.Sprintf("%v/%v/%v", splittedDate[2], splittedDate[1], splittedDate[0]) // dd/mm/yyyy

	hourMinSec := splittedTimeStr[1]
	splittedHourMinSec := strings.Split(hourMinSec, ":")
	newHourMin := fmt.Sprintf("%v:%v", splittedHourMinSec[0], splittedHourMinSec[1])

	str = fmt.Sprintf("%v %v", newDate, newHourMin)
	return
}
