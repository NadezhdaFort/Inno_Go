package main

import "fmt"

func main() {
	arr1 := []int{1, 2, 3, 4, 5}
	arr2 := []int{5, 4, 3, 2, 1}
	arr3 := []int{1, 3, 4, 5, 6}
	arr4 := []int{1, 2, 3}
	fmt.Println(IsEqualArrays(arr1, arr2))
	fmt.Println(IsEqualArrays(arr1, arr3))
	fmt.Println(IsEqualArrays(arr2, arr3))
	fmt.Println(IsEqualArrays(arr1, arr4))
}
func IsEqualArrays[T comparable](arr1, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	mapArr1 := make(map[T]int)
	for _, val := range arr1 {
		mapArr1[val]++
	}
	for _, val := range arr2 {
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
