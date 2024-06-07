package main

import (
	"fmt"
	"sync"
)

func mergeChannels[T any](ch1 <-chan T, ch2 <-chan T) <-chan T {
	var wg sync.WaitGroup
	merge := make(chan T)

	// Функция для чтения из канала и отправки в объединенный канал
	output := func(ch <-chan T) {
		defer wg.Done()
		for value := range ch {
			merge <- value
		}
	}

	wg.Add(2)
	go output(ch1)
	go output(ch2)

	// Закрытие объединенного канала после завершения всех операций слияния
	go func() {
		wg.Wait()
		close(merge)
	}()

	return merge
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
	}()
	go func() {
		defer close(ch2)
		for i := 20; i < 30; i++ {
			ch2 <- i
		}
	}()

	merge := mergeChannels(ch1, ch2)
	for n := range merge {
		fmt.Println(n)
	}
}
