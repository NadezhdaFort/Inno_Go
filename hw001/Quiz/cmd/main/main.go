package main

import (
	"bufio"
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

func main() {

	mix := flag.Bool("mix", false, "Mix up the questions")
	filename := flag.String("file", "problems.csv", "File name to read")
	flag.Parse()

	fileStatus := 0
	buf := 1024
	arr := make([][]string, 0, buf)
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Can't open file", err)
		return
	}
	defer file.Close()
	in := bufio.NewReader(os.Stdin)
	reader := csv.NewReader(file)
	var answer string
	countRight := 0
	countWrong := 0
	for {
		for i := 0; i < buf; i++ {
			record, err2 := reader.Read()
			if err2 == io.EOF {
				fileStatus = -1
				break
			}
			if err2 != nil {
				log.Fatal(err2)
			}
			arr = append(arr, record)
		}
		if *mix {
			rand.Shuffle(len(arr), func(i, j int) {
				arr[i], arr[j] = arr[j], arr[i]
			})
		}
		for i := 0; i < len(arr); i++ {
			fmt.Println(arr[i][0])
			fmt.Fscan(in, &answer)

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
			fmt.Printf("Number of wrong answers: %d", countWrong)
			return
		}
		arr = arr[:0]
	}
}
