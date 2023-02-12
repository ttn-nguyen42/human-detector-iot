package workers

import (
	"context"
)

/*
Add a task and let the program finds a routine for it to execute
*/

type TaskArguments struct {
	// The channel that hold the result of the task
	Res chan interface{}
}

type Task struct {
	// The task itself
	// Will be provided some arguments
	Do func(args *TaskArguments)

	// It can pass the result back into this channel
	// Must perform type casting in order to get the actual result
	Result chan interface{}
}

func NewTask(do func(args *TaskArguments)) Task {
	return Task{
		Do: do,
		Result: make(chan interface{}),
	}
}

func (t *Task) GetResult() interface{} {
	res := <-t.Result
	return res
}

type Worker struct {
	WorkerPool chan chan Task
	TaskQueue  chan Task
	Context    context.Context
	cancel     context.CancelFunc
}

// Pass in the worker pool that this worker belongs to
func NewWorker(workerPool chan chan Task) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		WorkerPool: workerPool,
		TaskQueue:  make(chan Task),
		cancel:     cancel,
		Context:    ctx,
	}
}

func (w *Worker) start() {
	go func(w *Worker) {

		for {
			// Register itself to the worker pool that it belongs to
			w.WorkerPool <- w.TaskQueue
			select {
			case task := <-w.TaskQueue:
				task.Do(&TaskArguments{
					Res: task.Result,
				})
				continue
			case <-w.Context.Done():

				return
			}
		}
	}(w)
}

func (w *Worker) Stop() {
	go func (w *Worker)  {
		w.cancel()
	}(w)
}
