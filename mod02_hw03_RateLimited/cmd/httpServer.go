// rate limit реализован с помощью 2-х map (requests map[string]int с счетчиком запросов и
// blockClients map[string]time.Time - списком блокировки, которая в качестве ключа принимает IP, а
// значение - время запроса + лимит ожидания(limitWaiting).
// При запросе проверяется map блокировки, если время ожидания прошло, то увеличивается счетчик запросов пока не
// достигнет limitRequests (в данном коде = 5), после чего IP перемещается в мапу блокировки со временем ожидания.
package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	limitRequests = 5
	limitWaiting  = 10 * time.Millisecond
)

var (
	requests     map[string]int
	blockClients map[string]time.Time
)

func init() {
	requests = make(map[string]int)
	blockClients = make(map[string]time.Time)
}

func getIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

func isBlocked(ip string) bool {
	if blockUntil, ok := blockClients[ip]; ok {
		if time.Now().Before(blockUntil) {
			return true
		}
		delete(blockClients, ip)
	}
	return false
}

func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ip := getIP(req.RemoteAddr)

		if isBlocked(ip) {
			http.Error(resp, "Blocked", http.StatusTooManyRequests)
			return
		}

		requests[ip]++
		if requests[ip] > limitRequests {
			blockClients[ip] = time.Now().Add(limitWaiting)
			requests[ip] = 0
			http.Error(resp, "Request limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(resp, req)
	}
}

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "hello\n")
}

func headers(resp http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(resp, "%v: %v\n", name, h)
		}
	}
}
func main() {
	http.HandleFunc("/hello", rateLimit(hello))
	http.HandleFunc("/headers", rateLimit(headers))

	http.ListenAndServe(":8080", nil)
}
