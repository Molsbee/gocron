package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse_EachMinute(t *testing.T) {
	// arrange
	now := time.Now()

	schedule := Parse("* * * * * *")
	start := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local)
	startOne := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 1, 0, 0, time.Local)
	startTwo := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 2, 0, 0, time.Local)
	startThree := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 3, 0, 0, time.Local)

	// act
	next := schedule.Next(start)
	nextOne := schedule.Next(startOne)
	nextTwo := schedule.Next(startTwo)
	nextThree := schedule.Next(startThree)

	// assert
	assert.True(t, next.Minute() == 0)
	assert.True(t, nextOne.Minute() == 1)
	assert.True(t, nextTwo.Minute() == 2)
	assert.True(t, nextThree.Minute() == 3)
}

func TestParse_MultipleMinutesProvided(t *testing.T) {
	// arrange
	now := time.Now()

	schedule := Parse("10,15 * * * * *")
	startTen := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 6, 0, 0, time.Local)
	startFifteen := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 13, 0, 0, time.Local)

	// act
	ten := schedule.Next(startTen)
	fifteen := schedule.Next(startFifteen)

	// assert
	assert.True(t, ten.Minute() == 10, fmt.Sprintf("Time: %s - minutes expected %d actual %d", ten, 10, ten.Minute()))
	assert.True(t, fifteen.Minute() == 15, fmt.Sprintf("Time: %s - minutes expected %d actual %d", fifteen, 15, fifteen.Minute()))
}

func TestParse_EveryMinuteOfFourthHour(t *testing.T) {
	// arrange
	now := time.Now()

	schedule := Parse("* 4 * * * *")
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	// act
	next := schedule.Next(start)

	// assert
	assert.Equal(t, 4, next.Hour()) // Executed on the 4th hour
}

func TestParse_EveryFifteenMinutesStartingAtZero(t *testing.T) {
	// arrange
	now := time.Now()

	schedule := Parse("*/15 * * * * *")
	startZero := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local)
	startFifteen := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 15, 0, 0, time.Local)
	startThirty := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 30, 0, 0, time.Local)
	startFourtyFive := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 45, 0, 0, time.Local)

	// act
	nextZero := schedule.Next(startZero)
	nextFifteen := schedule.Next(startFifteen)
	nextThirty := schedule.Next(startThirty)
	nextFourty := schedule.Next(startFourtyFive)

	// assert
	assert.True(t, nextZero.Before(startZero.Add(1*time.Minute)))
	assert.True(t, nextZero.Minute()%15 == 0)

	assert.True(t, nextFifteen.Before(startFifteen.Add(1*time.Minute)))
	assert.True(t, nextFifteen.Minute()%15 == 0)

	assert.True(t, nextThirty.Before(startThirty.Add(1*time.Minute)))
	assert.True(t, nextThirty.Minute()%15 == 0)

	assert.True(t, nextFourty.Before(startFourtyFive.Add(1*time.Minute)))
	assert.True(t, nextFourty.Minute()%15 == 0)
}
