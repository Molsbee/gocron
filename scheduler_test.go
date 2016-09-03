package gocron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduleTask_OneSecond(t *testing.T) {
	// arrange
	executionTime := make(chan time.Time, 2)

	// act
	scheduler := NewScheduler(1)
	scheduler.ScheduleSimpleTask().Second(1).Run(func() {
		t.Logf("Task is being executed %s", time.Now())
		executionTime <- time.Now()
	})

	t.Log("Starting Scheduler")
	stop := scheduler.Start()
	time.Sleep(2 * time.Second)
	stop <- true
	t.Log("Stopped Scheduler")

	// assert
	firstTime := <-executionTime
	assert.NotNil(t, firstTime)

	secondTime := <-executionTime
	assert.NotNil(t, secondTime)

	diff := secondTime.Sub(firstTime)
	t.Logf("Time difference %s", diff)
	assert.True(t, diff.Seconds() > 1)
	assert.True(t, diff.Seconds() < 2)
}

func TestScheduleTask_OneMinute(t *testing.T) {
	// arrange
	executionTime := make(chan time.Time, 2)

	// act
	scheduler := NewScheduler(1)
	scheduler.ScheduleSimpleTask().Minute(1).Run(func() {
		t.Logf("Task is being executed %s", time.Now())
		executionTime <- time.Now()
	})

	t.Log("Starting Scheduler")
	stop := scheduler.Start()
	time.Sleep(2 * time.Minute)
	stop <- true
	t.Log("Stopped Scheduler")

	// assert
	firstTime := <-executionTime
	assert.NotNil(t, firstTime)

	secondTime := <-executionTime
	assert.NotNil(t, secondTime)

	diff := secondTime.Sub(firstTime)
	t.Logf("Time difference %s", diff)
	assert.True(t, diff.Minutes() > 1)
	assert.True(t, diff.Minutes() < 2)
}

func TestScheduleTask_OneMinuteThirtySeconds(t *testing.T) {
	// arrange
	executionTime := make(chan time.Time, 2)

	// act
	scheduler := NewScheduler(1)
	scheduler.ScheduleSimpleTask().Minute(1).Second(30).Run(func() {
		t.Logf("Task is being executed %s", time.Now())
		executionTime <- time.Now()
	})

	t.Log("Starting Scheduler")
	stop := scheduler.Start()
	time.Sleep(2 * time.Minute)
	stop <- true
	t.Log("Stopped Scheduler")

	// assert
	firstTime := <-executionTime
	assert.NotNil(t, firstTime)

	secondTime := <-executionTime
	assert.NotNil(t, secondTime)

	diff := secondTime.Sub(firstTime)
	t.Logf("Time difference %s", diff)
	assert.True(t, diff.Minutes() > 1.3)
	assert.True(t, diff.Minutes() < 2)
}
