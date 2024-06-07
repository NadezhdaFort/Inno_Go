package main

import (
	"bufio"
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
func writeFile(filePath string, writeChan <-chan string, doneChan chan<- bool) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Невозможно открыть файл", err)
		doneChan <- true
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
			doneChan <- true
			return
		}
	}
	doneChan <- true
}

func main() {
	filePath := "output.txt"
	console := make(chan string)
	done := make(chan bool)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go readConsole(console)
	go writeFile(filePath, console, done)

	select {
	case <-signalChan:
		fmt.Println("Получен сигнал завершения работы")
	case <-done:
		fmt.Println("Программа завершена")
	}

	close(console)

	// ожидание записи в файл
	<-done
	os.Exit(0)
}
