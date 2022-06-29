package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	workersCount int
	jobs         chan Job
	results      chan map[string]int
	Done         chan struct{}
}

func NewWorkerPool(numberOfWorkers, numberOfJobs int) *WorkerPool {
	return &WorkerPool{
		workersCount: numberOfWorkers,
		jobs:         make(chan Job, numberOfJobs),
		results:      make(chan map[string]int, numberOfJobs),
		Done:         make(chan struct{}),
	}
}

func (wp *WorkerPool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)

		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

func (wp *WorkerPool) GenerateFrom(jobsBatch []Job) {
	for i, _ := range jobsBatch {
		wp.jobs <- jobsBatch[i]
	}

	close(wp.jobs)
}

func (wp WorkerPool) Results() <-chan map[string]int {
	return wp.results
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- map[string]int) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			// fan-in job execution multiplexing results into the results channel
			res, err := job.execute()
			if err != nil {
				return
			}

			results <- res

		case <-ctx.Done():
			return
		}
	}
}
