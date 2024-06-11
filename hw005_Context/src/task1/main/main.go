package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Функция чтения с консоли
func readConsole(inputChan chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите данные (для выхода из программы нажмите ctrl+c)")
	for scanner.Scan() {
		input := scanner.Text()
		inputChan <- input
	}
	if scanner.Err() != nil {
		fmt.Println("Ошибка ввода данных", scanner.Err())
		close(inputChan)
	}
}

// функция записи в файл
func writeFile(ctx context.Context, filePath string, writeChan <-chan string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Невозможно открыть файл", err)
		ctx.Done()
		return
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Println("Невозможно закрыть файл", err)
		}
	}(file)

	for text := range writeChan {
		if _, err := file.WriteString(text + "\n"); err != nil {
			fmt.Println("Невозможно записать в файл")
			ctx.Done()
			return
		}
	}
	ctx.Done()
}

func main() {
	filePath := "output.txt"
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	console := make(chan string)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go readConsole(console)
	go writeFile(ctx, filePath, console)

	select {
	case <-signalChan:
		fmt.Println(" Получен сигнал завершения работы")
	case <-ctx.Done():
		fmt.Println("Программа завершена")
	}

	close(console)

	// ожидание записи в файл
	cancel()
	os.Exit(0)
}
