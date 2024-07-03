package model

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Worker struct {
	cache        *Cache
	interval     time.Duration
	wg           sync.WaitGroup
	shutdownChan chan struct{}
}

func NewWorker(cache *Cache, interval time.Duration) *Worker {
	return &Worker{
		cache:        cache,
		interval:     interval,
		shutdownChan: make(chan struct{}),
	}
}

// Start запуск Worker-а, вызывается WorkerFactory
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.FlushCache()
			case <-w.shutdownChan:
				w.FlushCache()
				return
			}
		}
	}()
}

// Stop остановка Worker-а, вызавается WorkerFactory
func (w *Worker) Stop() {
	close(w.shutdownChan)
	w.wg.Wait()
}

// FlushCache получает map-у []Message из кэша, вызывает функцию записи в файл
func (w *Worker) FlushCache() {
	messages := w.cache.GetMessages()
	var wg sync.WaitGroup
	for fileName, message := range messages {
		wg.Add(1)
		go func(fileName string, message []Message) {
			defer wg.Done()
			w.WriteToFile(fileName, message)
		}(fileName, message)
	}
	wg.Wait()
}

// WriteToFile записывает данные из []Message в файл, если такого файла не существует, то создает его
func (w *Worker) WriteToFile(fileName string, msg []Message) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Ошибка открытия файла %s: %+v\n", fileName, err)
	}

	defer file.Close()

	for _, message := range msg {
		_, err := file.WriteString(fmt.Sprintf("%s: %s\n", message.Token, message.Data))
		if err != nil {
			fmt.Printf("Ошибка записи в файл: %s: %+v\n", fileName, err)
			return
		}
	}
}
