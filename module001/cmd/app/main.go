package main

import (
	. "Inno_Go/module001/cmd/app/model"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// WorkerPoolSize определяет количество воркеров
const WorkerPoolSize = 7

// TimeoutDuration определяет таймаут для каждого воркера
const TimeoutDuration = 1 * time.Second

// ValidToken валидирует токены
func ValidToken(token string) bool {
	if strings.HasPrefix(token, "validToken") {
		return true
	}
	return false
}

func main() {
	cache := NewCache()
	messageChannel := make(chan Message, 50)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	workerFactory := NewWorkerFactory(cache, TimeoutDuration, WorkerPoolSize)
	workerFactory.Start()

	// симуляция отправки сообщений
	go func() {
		tokens := []string{"validToken1", "validToken2", "validToken3", "validToken4", "validToken5", "invalidToken"}
		files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt", "file6.txt", "file7.txt", "file8.txt"}

		for i := 0; i < 100; i++ {
			token := tokens[i%len(tokens)]
			file := files[i%len(files)]
			messageChannel <- Message{Token: token,
				FileID: file,
				Data:   fmt.Sprintf("%s sent a message №%d to the file \"%s\"", token, i, file),
			}
		}
		close(messageChannel)
	}()

	// Добавление сообщений в кэш с проверкой токенов
	go func() {
		for msg := range messageChannel {
			if ValidToken(msg.Token) {
				cache.AddMessages(msg)
				// т.к. пришло сообщение, которое нужно записать в файл, сообщаем об этом монитору загрузок, чтобы он
				//через monitorLoad() создал еще один Worker, если их количество не превышает максимальное
				workerFactory.WorkersLoadChan <- struct{}{}
			} else {
				fmt.Printf("Invalid token: %s\n", msg.Token)
			}
		}
	}()

	<-signalChan
	workerFactory.Stop()
	fmt.Println(" Graceful Shutdown")
}
