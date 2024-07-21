// rate limit реализован с помощью limitRequests2 map[string]time.Time, которая в качестве ключа принимает IP, а
// значение - время запроса. При запросе проверяется время из map, соответствующее данному IP и сравнивается со
// временем запроса и если время запроса по временной шкале находится раньше времени, полученного из map, то данный
// запрос не обрабатывается.
// Если же запрос первый или время запроса по временной шкале находится позже времени, полученного из map, то в map
// сохраняется новое время с учетом limitWaiting2 для этого IP и запрос обрабатывается.
// Для удобства тестирования limitWaiting2 в данном коде  =  10 милисекундам.
package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	limitWaiting2 = 10 * time.Millisecond
)

var (
	limitRequests2 map[string]time.Time
)

func init() {
	limitRequests2 = make(map[string]time.Time)
}

func getIP2(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

func rateLimit2(next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := getIP2(req.RemoteAddr)
		timeRequest := time.Now()

		if t, ok := limitRequests2[ip]; ok {
			if timeRequest.Before(t) {
				http.Error(resp, "Blocked\n", http.StatusTooManyRequests)
				return
			}
		}

		limitRequests2[ip] = timeRequest.Add(limitWaiting2)
		next.ServeHTTP(resp, req)
	}
}

func hello2(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "hello\n")
}

func headers2(resp http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, header := range headers {
			fmt.Fprintf(resp, "%v: %v\n", name, header)
		}
	}
}

func main() {
	http.HandleFunc("/hello", rateLimit2(hello2))
	http.HandleFunc("/headers", rateLimit2(headers2))

	http.ListenAndServe(":8080", nil)
}
