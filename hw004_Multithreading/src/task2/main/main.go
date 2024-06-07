package main

import (
	"fmt"
	"math/big"
	"sync"
)

func sortNumbers(arr []int, prime chan<- int, composite chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(prime)
	defer close(composite)
	for _, number := range arr {
		if big.NewInt(int64(number)).ProbablyPrime(0) {
			prime <- number
		} else {
			composite <- number
		}
	}
}

func main() {
	var wg sync.WaitGroup
	arr := []int{3, 52, 13, 17, 44, 27, 85}
	slicePrimes := make([]int, 0, len(arr))
	sliceComposites := make([]int, 0, len(arr))

	primes := make(chan int)
	composites := make(chan int)

	// запуск горутины для разделения чисел
	wg.Add(1)
	go sortNumbers(arr, primes, composites, &wg)

	// горутина для сбора простых чисел
	go func() {
		for prime := range primes {
			slicePrimes = append(slicePrimes, prime)
		}
	}()

	// горутина для сбора составных чисел
	go func() {
		for composite := range composites {
			sliceComposites = append(sliceComposites, composite)
		}
	}()

	wg.Wait()

	fmt.Println("Primes:", slicePrimes)
	fmt.Println("Composites:", sliceComposites)
}
