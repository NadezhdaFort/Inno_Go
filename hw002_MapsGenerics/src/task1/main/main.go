package main

import (
	"fmt"
	"sort"
)

func findIntersection(slices ...[]int) []int {
	if len(slices) == 0 {
		fmt.Print("Не введенно ни одного слайса ")
		return nil
	}
	if len(slices[0]) == 0 {
		return []int{}
	}
	intersection := make(map[int]struct{}, len(slices[0]))
	for _, num := range slices[0] {
		intersection[num] = struct{}{}
	}
	for _, slice := range slices[1:] {
		tempMap := make(map[int]struct{}, len(intersection))
		for _, num := range slice {
			if _, ok := intersection[num]; ok {
				tempMap[num] = struct{}{}
			}
		}
		intersection = tempMap
	}
	result := make([]int, 0, len(intersection))
	for n, _ := range intersection {
		result = append(result, n)
	}
	sort.Ints(result)
	return result
}

func main() {
	arr := findIntersection()
	if arr == nil {
		fmt.Println()
	}
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 2}))
	fmt.Println(findIntersection([]int{1, 2, 3, 4}, []int{3, 2}, []int{}))
}
