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
func (numbers Numbers[T]) Equal(numbers2 Numbers[T]) bool {
	if len(numbers) != len(numbers2) {
		return false
	}
	mapArr1 := make(map[T]int)
	for _, val := range numbers {
		mapArr1[val]++
	}
	for _, val := range numbers2 {
		if _, ok := mapArr1[val]; !ok {
			return false
		}
		mapArr1[val]--
		if mapArr1[val] == 0 {
			delete(mapArr1, val)
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
