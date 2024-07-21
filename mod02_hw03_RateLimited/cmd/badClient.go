package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	for i := 0; i < 500; i++ {
		resp, err := http.Get("http://localhost:8080/headers")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Request № ", i)
		fmt.Println("Response status: ", resp.Status)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		log.Println("body ", string(body))
	}
}
