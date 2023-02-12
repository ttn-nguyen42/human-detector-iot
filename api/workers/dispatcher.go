package workers

import "iot_api/utils"

// For sending the task into
var taskQueue chan Task

func GetTaskQueue() chan Task {
	return taskQueue
}

func AddTask(task Task) {
	go func() {
		GetTaskQueue() <- task
	}()
}

// Handles the distribution of task to the workers
type Dispatcher struct {
	WorkerPool chan chan Task
	maxWorkers int
	cancel     []func()
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Task, maxWorkers)
	return &Dispatcher{
		WorkerPool: pool,
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Breath() {
	utils.Repeat(d.maxWorkers, func(at int) {
		// Create worker that belongs to this dispatcher
		worker := NewWorker(d.WorkerPool)
		worker.start()
		d.cancel = append(d.cancel, worker.Stop)
	})
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for task := range GetTaskQueue() {
		go func(task Task) {
			// Finds an available worker
			avail := <-d.WorkerPool
			// Assign a task to it
			avail <- task
		}(task)
	}
}

func (d *Dispatcher) Kill() {
	utils.ForEach(&d.cancel, func(with func()) {
		go func(cancel func()) {
			cancel()
		}(with)
	})
}
