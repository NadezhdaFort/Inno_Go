package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func stringPreparation(str string) string {
	s := strings.ToLower(str)
	s = strings.TrimSpace(s)
	return s
}
func loadData(reader *csv.Reader, arr [][]string, buf int) [][]string {
	for i := 0; i < buf; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			fileStatus = -1
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, record)
	}
	return arr
}

var fileStatus = 0

const buf = 1024

func main() {

	mix := flag.Bool("mix", false, "Mix up the questions")
	filename := flag.String("file", "problems.csv", "File name to read")
	flag.Parse()

	arr := make([][]string, 0, buf)
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal("Can't open file", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Can't exit file", err)
			return
		}
	}(file)

	reader := csv.NewReader(file)
	var answer string
	countRight := 0
	countWrong := 0
	for {
		arr := loadData(reader, arr, buf)
		if *mix {
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
		}
		if len(arr) == 0 {
			fmt.Println("Используемый файл не содержит данных или поврежден")
			return
		}
		for i := 0; i < len(arr); i++ {
			fmt.Println(arr[i][0])
			_, err := fmt.Scanln(&answer)
			if err != nil {
				return
			}

			a := stringPreparation(answer)
			b := stringPreparation(arr[i][1])

			if a == b {
				countRight++
			} else {
				countWrong++
			}
		}
		if fileStatus == -1 {
			fmt.Printf("Number of right answers: %d\n", countRight)
			fmt.Printf("Number of wrong answers: %d\n", countWrong)
			return
		}
		arr = arr[:0]
	}
}
