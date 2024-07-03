// WorkerFactory управляет количеством Worker-ов
package model

import (
	"sync"
	"time"
)

type WorkerFactory struct {
	mu              sync.Mutex
	workers         []*Worker
	cache           *Cache
	interval        time.Duration
	wg              sync.WaitGroup
	WorkersLoadChan chan struct{}
	shutdownChan    chan struct{}
	maxCountWorkers int
}

func NewWorkerFactory(cache *Cache, interval time.Duration, maxCountWorkers int) *WorkerFactory {
	return &WorkerFactory{
		cache:           cache,
		interval:        interval,
		WorkersLoadChan: make(chan struct{}, maxCountWorkers),
		shutdownChan:    make(chan struct{}),
		maxCountWorkers: maxCountWorkers,
	}
}

// Start запускает создание Worker-а и отслеживание нагрузки
func (wf *WorkerFactory) Start() {
	// запуск первого Worker-а
	wf.AddWorker()

	wf.wg.Add(1)
	go wf.monitorLoad()
}

// Stop вызывается при graceful shutdown и вызывает метод Stop() у всех созданных Worker-ов
func (wf *WorkerFactory) Stop() {
	close(wf.shutdownChan)
	defer close(wf.WorkersLoadChan)

	wf.mu.Lock()
	for _, worker := range wf.workers {
		worker.Stop()
	}
	wf.mu.Unlock()
	wf.wg.Wait()
}

// monitorLoad вызывается если есть данные в кэше, но все доступные Worker-ы заняты, и если количество созданных
// Worker-ов не превышает максимально допустимое значение, то создается еще один Worker
func (wf *WorkerFactory) monitorLoad() {
	defer wf.wg.Done()
	for {
		select {
		case <-wf.WorkersLoadChan:
			if len(wf.workers) < wf.maxCountWorkers {
				wf.AddWorker()
			}

		case <-wf.shutdownChan:
			return
		}
	}
}

// AddWorker создает новый Worker и добавляет его в слайс Worker-ов в WorkerFactory и запускает
func (wf *WorkerFactory) AddWorker() {
	worker := NewWorker(wf.cache, wf.interval)
	wf.mu.Lock()
	wf.workers = append(wf.workers, worker)
	wf.mu.Unlock()

	worker.Start()
}
