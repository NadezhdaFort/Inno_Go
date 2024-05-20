package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	create := flag.String("create", "", "Create file")
	read := flag.String("read", "", "Read file")
	remove := flag.String("remove", "", "Remove file")
	flag.Parse()
	if *create != "" {
		createFile(*create)
	}
	if *remove != "" {
		removeFile(*remove)
	}
	if *read != "" {
		readFile(*read)
	}
}
func createFile(s string) {
	if _, er := os.Stat(s); os.IsNotExist(er) {
		file, err := os.Create(s)
		if err != nil {
			fmt.Println("Ошибка создания файла:", err)
			return
		}
		defer file.Close()
		fmt.Println("Файл", s, "успешно создан")
	} else {
		fmt.Println("Файл", s, "уже существует")
	}
}
func readFile(s string) {
	file, err := os.Open(s)
	if err != nil {
		fmt.Println("Невозможно открыть файл", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла", err)
		return
	}
	fmt.Println("Файл", s, "успешно прочитан")
}
func removeFile(s string) {
	err := os.Remove(s)
	if err != nil {
		fmt.Println("Ошибка удаления файла:", err)
		return
	}
	fmt.Println("Файл", s, "успешно удален")
}
