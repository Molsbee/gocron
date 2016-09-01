package gocron

import "time"

// Scheduler contains multiple tasks for execution and manages those tasks
// execution
type Scheduler interface {
	Schedule() *Task
	RunPending()
	Start() chan bool
}

// MaxTasks is the current max size of a task list that an array can administer
const MaxTasks = 100

type scheduler struct {
	tasks [MaxTasks]*Task
	size  int
}

// NewScheduler -
func NewScheduler() Scheduler {
	return &scheduler{
		tasks: [MaxTasks]*Task{},
		size:  0,
	}
}

func (s *scheduler) Schedule() *Task {
	task := NewTask()
	s.tasks[s.size] = task
	s.size++
	return task
}

func (s *scheduler) RunPending() {
	// TODO: Should delegate work to multiple works
	for i := 0; i < s.size; i++ {
		if s.tasks[i].shouldRun() {
			s.tasks[i].run()
		}
	}
}

func (s *scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.RunPending()
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}
