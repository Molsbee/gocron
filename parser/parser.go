package parser

import (
	"log"
	"math"
	"strconv"
	"strings"
)

type boundary struct {
	min   uint
	max   uint
	names map[string]uint
}

var (
	minutes    = boundary{min: 0, max: 59, names: nil}
	hours      = boundary{min: 0, max: 23, names: nil}
	dayOfMonth = boundary{min: 1, max: 31, names: nil}
	months     = boundary{
		min: 1,
		max: 12,
		names: map[string]uint{
			"jan": 1,
			"feb": 2,
			"mar": 3,
			"apr": 4,
			"may": 5,
			"jun": 6,
			"jul": 7,
			"aug": 8,
			"sep": 9,
			"oct": 10,
			"nov": 11,
			"dec": 12,
		},
	}
	dayOfWeek = boundary{
		min: 1,
		max: 7,
		names: map[string]uint{
			"sun": 1,
			"mon": 2,
			"tue": 3,
			"wed": 4,
			"thu": 5,
			"fri": 6,
			"sat": 7,
		},
	}
	year = boundary{min: 1900, max: 3000, names: nil}
)

// Parse converts a classic cron formatted string into a schedule which
// provides mechanism to determine the next time the task should be executed.
// * * * * * *
// | | | | | |
// | | | | | +-- Year              (range: 1900-3000)
// | | | | +---- Day of the Week   (range: 1-7, 1 standing for Monday)
// | | | +------ Month of the Year (range: 1-12)
// | | +-------- Day of the Month  (range: 1-31)
// | +---------- Hour              (range: 0-23)
// +------------ Minute            (range: 0-59)
func Parse(cron string) *Schedule {
	fields := strings.Fields(cron)
	if len(fields) != 6 {
		log.Panic("Incorrect schedule provided")
	}

	return &Schedule{
		Minute:     getField(fields[0], minutes),
		Hour:       getField(fields[1], hours),
		DayOfMonth: getField(fields[2], dayOfMonth),
		Month:      getField(fields[3], months),
		DayOfWeek:  getField(fields[4], dayOfWeek),
		Year:       getField(fields[5], year),
	}
}

func getField(v string, b boundary) uint64 {
	values := strings.FieldsFunc(v, func(r rune) bool {
		return r == ','
	})

	if len(values) != 1 {
		var bits uint64
		for _, value := range values {
			bits |= getField(value, b)
		}
		return bits
	}

	var start, end, step uint
	stepValues := strings.Split(v, "/")
	rangeValues := strings.Split(stepValues[0], "-")

	if rangeValues[0] == "*" || rangeValues[0] == "?" {
		start = b.min
		end = b.max
	} else {
		start = parseField(rangeValues[0], b.names)
		switch len(rangeValues) {
		case 1:
			end = start
		case 2:
			end = parseField(rangeValues[1], b.names)
		default:
			log.Panicf("Too many hyphens provided: %s", v)
		}

	}

	if start < b.min {
		log.Panicf("Value (%d) provided is below the minimum allowed value (%d) for field", start, b.min)
	}
	if end > b.max {
		log.Panicf("Value (%d) provided is above the maximum allowed value (%d) for field", end, b.max)
	}
	if start > end {
		log.Panicf("Beginning of range (%d) is beyond the end range (%d)", start, end)
	}

	switch len(stepValues) {
	case 1:
		step = 1
	case 2:
		step = parseField(stepValues[1], nil)
	default:
		log.Panicf("To many slashes provided: %s", v)
	}

	return getBits(start, end, step)
}

func parseField(field string, names map[string]uint) uint {
	if names != nil {
		if value, ok := names[strings.ToLower(field)]; ok {
			return value
		}
	}
	value, err := strconv.Atoi(field)
	if err != nil || value < 0 {
		log.Panicf("Failed to parse field or incorrect value provided %s", field)
	}

	return uint(value)
}

func getBits(min, max, step uint) uint64 {
	var bits uint64

	// If step is 1, use shifts.
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	// Else, use a simple loop.
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}
