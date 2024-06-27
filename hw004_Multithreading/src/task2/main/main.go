package main

import (
	"fmt"
	"math/big"
	"sync"
)

func sortNumbers(arr []int, wg *sync.WaitGroup) (p <-chan int, c <-chan int) {
	primes := make(chan int)
	composites := make(chan int)

	go func() {
		defer wg.Done()
		defer close(primes)
		defer close(composites)

		for _, number := range arr {
			if big.NewInt(int64(number)).ProbablyPrime(0) {
				primes <- number
			} else {
				composites <- number
			}
		}
	}()
	return primes, composites
}

func main() {
	var wg sync.WaitGroup
	arr := []int{1, 2, 3, 4, 5, 6, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	slicePrimes := make([]int, 0, len(arr))
	sliceComposites := make([]int, 0, len(arr))

	// запуск горутины для разделения чисел
	wg.Add(3)
	primes, composites := sortNumbers(arr, &wg)

	// горутина для сбора простых чисел
	go func() {
		defer wg.Done()
		for prime := range primes {
			slicePrimes = append(slicePrimes, prime)
		}
	}()

	// горутина для сбора составных чисел
	go func() {
		defer wg.Done()
		for composite := range composites {
			sliceComposites = append(sliceComposites, composite)
		}
	}()

	wg.Wait()

	fmt.Println("Primes:", slicePrimes)
	fmt.Println("Composites:", sliceComposites)
}
