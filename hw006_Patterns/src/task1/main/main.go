package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// WorkerPoolSize определяет количество воркеров
const WorkerPoolSize = 3

// TimeoutDuration определяет таймаут для каждого воркера
const TimeoutDuration = 2 * time.Second

// URLs перечень файлов для загрузки
var URLs = []string{
	"https://img.freepik.com/free-photo/high-angle-delicious-pancakes-arrangement_23-2150265090.jpg?t=st=1718634824~exp=17",
	"https://ru.wikipedia.org/wiki/%D0%A4%D0%B0%D0%B9%D0%BB:Text-txt.svg",
	"https://upload.wikimedia.org/wikipedia/commons/5/54/Panda_Cub_%284274178112%29.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTRMYE4HcgwCeTPZhLe_J86TinY1IGTqsjr4LMSZE9Pwz82KNmSJ4Q1JCYEeyypVS9mj-8&usqp=CAU",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQEJQknYlIFpDl2fLdbIbmwn9fJ3E9g7qslUg&s",
	"https://www.cleverfiles.com/howto/wp-content/uploads/2018/03/minion.jpg",
	"https://desano.ru/uploads/catalog/309/NS-10051-1.jpg",
}

func downloadFile(url string, wg *sync.WaitGroup, results chan<- string) error {
	defer wg.Done()

	client := http.Client{
		Timeout: TimeoutDuration,
	}

	resp, err := client.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Failed to download %s: %v", url, err)
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing resp.Body")
		}
	}(resp.Body)

	out, err := os.Create("downloaded_" + extractFileName(url))
	if err != nil {
		results <- fmt.Sprintf("Failed to create file for %s: %v", url, err)
		return err
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Failed to write file for %s: %v", url, err)
		return err
	}

	results <- fmt.Sprintf("Successfully downloaded %s", url)
	return nil
}

// extractFileName функция для извлечения имени файла из URL
func extractFileName(url string) string {
	segments := strings.Split(url, "/")
	return segments[len(segments)-1]
}

// worker функция, выполняющая загрузку файлов
func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for url := range jobs {
		fmt.Printf("Worker %d started downloading %s\n", id, url)
		if err := downloadFile(url, wg, results); err != nil {
			fmt.Println("File download error", err)
			continue
		}
		fmt.Printf("Worker %d finished downloading %s\n", id, url)
	}
}

func main() {
	var wg sync.WaitGroup

	// Каналы для задач и результатов
	jobs := make(chan string, 3)
	results := make(chan string, 3)

	// Запуск воркеров
	for w := 1; w <= WorkerPoolSize; w++ {
		go worker(w, jobs, results, &wg)
	}

	// Добавление задач в канал jobs
	wg.Add(len(URLs))
	for _, url := range URLs {
		jobs <- url
	}
	close(jobs)

	// Ожидание завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Вывод результатов
	for result := range results {
		fmt.Println(result)
	}
}
