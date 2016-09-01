package gocron

import "time"

// Scheduler contains multiple tasks for execution and manages those tasks
// execution
type Scheduler interface {
	Schedule() *Task
	Start() chan bool
}

type scheduler struct {
	tasks       []*Task
	workChannel chan *Task
}

// NewScheduler creates a new instance of a Scheduler which consists of multiple
// task for execution.  Along with a pool of workers defined for processing those
// tasks.  The pool is defined through this constractor and it's recomended to
// minimize the amount of time intensive tasks as they will cause the workers
// to delay execution.  Size tasks and worker pool appropriately
func NewScheduler(poolSize int) Scheduler {
	scheduler := &scheduler{
		tasks:       []*Task{},
		workChannel: make(chan *Task, 10),
	}

	for i := 0; i <= poolSize; i++ {
		go scheduler.processTask()
	}

	return scheduler
}

// Schedule is used as a builder pattern to return a task that can then be
// added to the Scheduler through a fluent style api
func (s *scheduler) Schedule() *Task {
	task := NewTask()
	s.tasks = append(s.tasks, task)
	return task
}

func (s *scheduler) runPending() {
	// TODO: Should delegate work to multiple works
	for i := 0; i < len(s.tasks); i++ {
		if s.tasks[i].shouldRun() {
			s.workChannel <- s.tasks[i]
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
				s.runPending()
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

func (s *scheduler) processTask() {
	for task := range s.workChannel {
		task.run()
	}
}
