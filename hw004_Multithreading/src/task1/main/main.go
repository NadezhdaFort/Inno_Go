package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func readConsole() <-chan string {
	inputChan := make(chan string)

	go func() {
		defer close(inputChan)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Для завершения работы введите \"Ctrl + C\"")

		for scanner.Scan() {
			text := scanner.Text()
			inputChan <- text
		}
		if scanner.Err() != nil {
			fmt.Println("Ошибка ввода данных", scanner.Err())
			return
		}
	}()
	return inputChan
}

func writeFile(fileName string, inputChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Ошибка открытия файла", err)
		return
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Println("Невозможно закрыть файл", err)
		}
	}(file)

	for text := range inputChan {
		if _, err := file.WriteString(text + "\n"); err != nil {
			fmt.Println("Ошибка записи в файл", err)
			return
		}
	}
}

func main() {
	fileName1 := "output.txt"
	var wg sync.WaitGroup
	done := make(chan struct{})

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println(" Введен сигнал завершения работы")
		done <- struct{}{}
	}()

	readChan := readConsole()

	wg.Add(1)
	go writeFile(fileName1, readChan, &wg)

	go func() {
		wg.Wait()
	}()
	<-done
	close(done)
	fmt.Println("Программа завершена")
	os.Exit(0)
}
