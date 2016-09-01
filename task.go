package gocron

import (
	"log"
	"time"
)

// Task struct defining a operation that needs to be execute by gocron
type Task struct {
	task   Runnable
	hour   int
	minute int
	second int

	lastRun time.Time
	nextRun time.Time
	period  time.Duration
}

// NewTask -
func NewTask() *Task {
	return &Task{}
}

// Hour every x number of hours the cron task will be executed
func (t *Task) Hour(h int) *Task {
	if h < 0 || h > 25 {
		log.Panic("Hour time interval outside of bounds [0 - 24]")
	}
	t.hour = h
	return t
}

// Minute every x number of minutes the cron task will be executed
func (t *Task) Minute(m int) *Task {
	if m < 0 || m > 60 {
		log.Panic("Minute time interval outside of bounds [0 - 60]")
	}
	t.minute = m
	return t
}

// Second every x number of seconds the cron task will be executed
func (t *Task) Second(s int) *Task {
	if s < 0 || s > 60 {
		log.Panic("Second time interval outside of bounds [0-60]")
	}
	t.second = s
	return t
}

// Runnable contract for a function that is runnable by cron scheduler
type Runnable interface {
	Run()
}

// RunnableFunc type is an adapter to allow the use of a ordinary functions
// as argument to cron Run.
type RunnableFunc func()

// Run method executes function with no arguments
func (r RunnableFunc) Run() {
	r()
}

// Run specifies the function that should be called every time the task Scheduler
// executes.  Currently the only supported function will have to implement the
// runnable contract which is a no argument function with no return value.
// Most tasks being ran on cron won't have anyone to parse results so no return
// value is expected.  Also a function that is runnable should be self contained
// at this point should alow avoiding reflection
func (t *Task) Run(fn func()) {
	t.task = RunnableFunc(fn)
	t.scheduleNextExecution()
}

func (t *Task) run() {
	t.task.Run()
	t.lastRun = time.Now()
	t.scheduleNextExecution()
}

func (t *Task) shouldRun() bool {
	return time.Now().After(t.nextRun)
}

func (t *Task) scheduleNextExecution() {
	if t.lastRun == time.Unix(0, 0) {
		t.lastRun = time.Now()
	}

	if t.period == 0 {
		second := time.Duration(t.second)
		minute := time.Duration(t.minute * 60)
		hour := time.Duration(t.hour * 60 * 60)
		t.period = second + minute + hour
	}
	t.nextRun = t.lastRun.Add(t.period * time.Second)
}
