package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// WorkerPoolSize определяет количество воркеров
const WorkerPoolSize = 5

// TimeoutDuration определяет таймаут для каждого воркера
const TimeoutDuration = 2 * time.Second

// URLs перечень файлов для загрузки
var URLs = []string{
	"https://img.freepik.com/free-photo/high-angle-delicious-pancakes-arrangement_23-2150265090.jpg?t=st=1718634824~exp=17",
	"https://img.freepik.com/free-psd/3d-rendering-of-ui-icon_23-2149182289.jpg?w=1060&t=st=1718636284~exp=1718636884~hmac=fd18c9f95a809901e907b258f8867ec766cab6341788f191aeab3e11e0ee81bf",
	"https://img.freepik.com/free-photo/3d-cartoon-view-lawyer-briefcase_23-2151419584.jpg?t=st=1718636358~exp=1718639958~hmac=071dbb3dcc61a91eb78f2a6ed1f2a10821eb6b19e2f9c9b091ae7e1d0b3818dc&w=826",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTRMYE4HcgwCeTPZhLe_J86TinY1IGTqsjr4LMSZE9Pwz82KNmSJ4Q1JCYEeyypVS9mj-8&usqp=CAU",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQEJQknYlIFpDl2fLdbIbmwn9fJ3E9g7qslUg&s",
	"https://stickershop.line-scdn.net/stickershop/v1/product/8140588/LINEStorePC/main.png?v=1",
	"https://desano.ru/uploads/catalog/309/NS-10051-1.jpg",
}

// downloadFile функция для загрузки файла по URL
func downloadFile(url string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	client := http.Client{
		Timeout: TimeoutDuration,
	}

	resp, err := client.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Failed to download %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Failed to ride file %s: %v", url, err)
		return
	}
	results <- fmt.Sprintf("Successfully downloaded %d bytes", len(body))
}

// worker функция, выполняющая загрузку файлов
func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for url := range jobs {
		fmt.Printf("Worker %d started downloading %s\n", id, url)
		downloadFile(url, wg, results)
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
	for _, url := range URLs {
		wg.Add(1)
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
