package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func readConsole(ctx context.Context) <-chan string {
	inputChan := make(chan string)

	go func() {
		defer close(inputChan)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Для завершения работы введите \"Ctrl + C\"")

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				text := scanner.Text()
				inputChan <- text
			}
		}
		if scanner.Err() != nil {
			fmt.Println("Ошибка ввода данных", scanner.Err())
			return
		}
	}()
	return inputChan
}

func writeFile(ctx context.Context, fileName string, inputChan <-chan string, wg *sync.WaitGroup) {
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

	for {
		select {
		case text, ok := <-inputChan:
			if !ok {
				return
			}
			if _, err := file.WriteString(text + "\n"); err != nil {
				fmt.Println("Ошибка записи в файл", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	fileName1 := "output.txt"
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println(" Введен сигнал завершения работы")
		cancel()
	}()

	readChan := readConsole(ctx)

	wg.Add(1)
	go writeFile(ctx, fileName1, readChan, &wg)

	wg.Wait()

	fmt.Println("Программа завершена")
	os.Exit(0)
}
