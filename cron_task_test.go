package gocron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCronTask_ShouldRun(t *testing.T) {
	task := NewCronTask("* * * * * *", func() {
		t.Log("Hello")
	})

	shouldRun := task.shouldRun()
	assert.False(t, shouldRun)

	time.Sleep(1 * time.Second)

	shouldRun = task.shouldRun()
	assert.True(t, shouldRun)
}

func TestCronTask_Run(t *testing.T) {
	// arrange
	called := false
	task := NewCronTask("* * * * * *", func() {
		called = true
	})

	// act
	task.run()

	// assert
	assert.True(t, called)
}
