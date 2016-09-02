package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TimeData struct {
	cronSchedule string
	validation   func(t *testing.T, current, next time.Time)
	startTime    []time.Time
}

func executeTimeData(t *testing.T, timeData TimeData) {
	schedule := Parse(timeData.cronSchedule)
	for _, time := range timeData.startTime {
		next := schedule.Next(time)
		t.Logf("current (%s) next (%s)", time, next)
		timeData.validation(t, time, next)
	}
}

var now = time.Now()

func TestParse_EveryMinute(t *testing.T) {
	everyMinute := TimeData{
		cronSchedule: "* * * * * *",
		validation: func(t *testing.T, current, next time.Time) {
			assert.Equal(t, current.Year(), next.Year())
			assert.Equal(t, current.Hour(), next.Hour())
			assert.Equal(t, current.Minute(), next.Minute())
		},
		startTime: []time.Time{
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 5, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 12, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 25, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 45, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 59, 0, 0, time.Local),
		},
	}

	executeTimeData(t, everyMinute)
}

func TestParse_MultipleMinuteList(t *testing.T) {
	minuteList := TimeData{
		cronSchedule: "25,45 * * * * *",
		validation: func(t *testing.T, current, next time.Time) {
			assert.Equal(t, current.Year(), next.Year())
			assert.Equal(t, current.Month(), next.Month())
			assert.Equal(t, current.Day(), next.Day())

			if current.Minute() < 25 || current.Minute() > 45 {
				assert.True(t, next.Minute() == 25)
			} else {
				assert.True(t, next.Minute() == 45)
			}
		},
		startTime: []time.Time{
			time.Date(now.Year(), now.Month(), now.Day(), 0, 10, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 0, 30, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 0, 50, 0, 0, time.Local),
		},
	}

	executeTimeData(t, minuteList)
}

func TestParse_MinuteStepValue(t *testing.T) {
	minuteStep := TimeData{
		cronSchedule: "*/15 * * * * *",
		validation: func(t *testing.T, current, next time.Time) {
			assert.Equal(t, current.Year(), next.Year())
			assert.Equal(t, current.Month(), next.Month())
			assert.Equal(t, current.Day(), next.Day())
			assert.True(t, next.Minute()%15 == 0)
		},
		startTime: []time.Time{
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 5, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 20, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 35, 0, 0, time.Local),
		},
	}

	executeTimeData(t, minuteStep)
}

func TestParse_HourList(t *testing.T) {
	var hourList = TimeData{
		cronSchedule: "* 4,8 * * * *",
		validation: func(t *testing.T, current, next time.Time) {
			assert.Equal(t, current.Year(), next.Year())
			assert.Equal(t, current.Month(), next.Month())
			if current.Hour() < 4 || current.Hour() > 8 {
				day := current.Day()
				if current.Hour() > 8 {
					day++
				}
				assert.Equal(t, day, next.Day())
				assert.Equal(t, 4, next.Hour())
			} else {
				assert.Equal(t, current.Day(), next.Day())
				assert.Equal(t, 8, next.Hour())
			}
		},
		startTime: []time.Time{
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.Local),
			time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.Local),
		},
	}

	executeTimeData(t, hourList)
}

func TestParse_MonthFieldFebruary(t *testing.T) {
	var month = TimeData{
		cronSchedule: "* * * 2 * *",
		validation: func(t *testing.T, current, next time.Time) {
			assert.True(t, next.Month() == 2)

			if current.Month() != 2 {
				if current.Month() > 2 {
					assert.Equal(t, current.Year()+1, next.Year())
				} else {
					assert.Equal(t, current.Year(), next.Year())
				}

				assert.Equal(t, 1, next.Day())
				assert.Equal(t, 0, next.Hour())
				assert.Equal(t, 0, next.Minute())
			} else {
				assert.Equal(t, current.Year(), next.Year())
				assert.Equal(t, current.Hour(), next.Hour())
				assert.Equal(t, current.Minute(), next.Minute())
			}
		},
		startTime: []time.Time{
			time.Date(now.Year(), 2, 1, 0, 0, 0, 0, time.Local),
			time.Date(now.Year(), 2, 1, 0, 1, 0, 0, time.Local),
			time.Date(now.Year(), 2, 1, 0, 2, 0, 0, time.Local),
			time.Date(now.Year(), 1, 2, 0, 0, 0, 0, time.Local), // This should be Feb | day 1 | hour 0 | minute 0
			time.Date(now.Year(), 8, now.Day(), 0, 0, 0, 0, time.Local),
		},
	}

	executeTimeData(t, month)
}
