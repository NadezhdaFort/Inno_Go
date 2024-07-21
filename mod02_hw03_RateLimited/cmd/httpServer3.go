// rate limit реализован с помощью time.Ticker-а (в функции init()), который отправляет каждые 10 милисекунд
// (для удобства тестирования) в канал requests3 текущую временную метку.
// Функция rateLimit() является оберткой функций обработчиков HTTP. Если из канала приходят данные, то это означает,
// что время ожидания закончено и запрос можно обработать.
package main

import (
	"fmt"
	"net/http"
	"time"
)

var (
	requests3 chan time.Time
	ticker    *time.Ticker
)

func init() {
	requests3 = make(chan time.Time, 1)
	ticker = time.NewTicker(10 * time.Millisecond)
	go func() {
		for range ticker.C {
			requests3 <- time.Now()
		}
	}()
}

func rateLimit3(next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		select {
		case <-requests3:
			next.ServeHTTP(resp, req)
		default:
			http.Error(resp, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	}
}

func hello3(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "hello\n")
}

func headers3(resp http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, header := range headers {
			fmt.Fprintf(resp, "%v: %v\n", name, header)
		}
	}
}

func main() {
	http.HandleFunc("/hello", rateLimit3(hello3))
	http.HandleFunc("/headers", rateLimit3(headers3))

	http.ListenAndServe(":8080", nil)
}
