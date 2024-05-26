package main

import (
	"errors"
)

type Numbers[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64] []T

func (numbers Numbers[T]) Sum() T {
	var sum T
	for _, n := range numbers {
		sum += n
	}
	return sum
}
func (numbers Numbers[T]) Multiply() T {
	var multi T = 1
	for _, n := range numbers {
		multi *= n
	}
	return multi
}
func (numbers Numbers[T]) Equal(arr Numbers[T]) bool {
	if len(numbers) != len(arr) {
		return false
	}
	for i, val := range numbers {
		if val != arr[i] {
			return false
		}
	}
	return true
}
func (numbers Numbers[T]) Find(t T) (int, error) {
	for i, val := range numbers {
		if val == t {
			return i, nil
		}
	}
	return -1, errors.New("not found")
}
func (numbers *Numbers[T]) DeleteByValue(t T) (bool, error) {
	for i, val := range *numbers {
		if val == t {
			*numbers = append((*numbers)[:i], (*numbers)[i+1:]...)
			return true, nil
		}
	}
	return false, errors.New("wrong value")
}
func (numbers *Numbers[T]) DeleteByIndex(i int) (bool, error) {
	if i < 0 || i >= len(*numbers) {
		return false, errors.New("wrong index")
	}
	*numbers = append((*numbers)[:i], (*numbers)[i+1:]...)
	return true, nil
}
func main() {
}
