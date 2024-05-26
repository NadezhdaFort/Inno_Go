package main

import (
	"fmt"
	"sort"
)

func findIntersection(slices ...[]int) []int {
	intersection := make(map[int]struct{})
	for _, num := range slices[0] {
		intersection[num] = struct{}{}
	}
	for _, slice := range slices[1:] {
		tempArr := make(map[int]struct{})
		for _, num := range slice {
			if _, ok := intersection[num]; ok {
				tempArr[num] = struct{}{}
			}
		}
		intersection = tempArr
	}
	result := []int{}
	for n, _ := range intersection {
		result = append(result, n)
	}
	sort.Ints(result)
	return result
}

func main() {
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}, []int{}))
}
