package main

import (
	"errors"
)

func EvalSequence(matrix [][]int, userAnswer []int) (int, error) {

	if err := validation(matrix, userAnswer); err != nil {
		return 0, err
	}

	maxGrade := calMaxGrade(matrix)
	userGrade := calcUserGrade(matrix, userAnswer)

	percent := userGrade * 100 / maxGrade

	return percent, nil
}

func calMaxGrade(matrix [][]int) int {
	type Path struct {
		nodes  []int
		weight int
	}

	n := len(matrix)
	maxWeight := 0
	visited := make(map[int]struct{})

	var bfs func(start int)
	bfs = func(start int) {
		queue := []Path{{nodes: []int{start}, weight: 0}}

		for len(queue) > 0 {
			currentPath := queue[0]
			lastNode := currentPath.nodes[len(currentPath.nodes)-1]
			queue = queue[1:]

			for _, node := range currentPath.nodes {
				visited[node] = struct{}{}
			}

			needsToBeQueued := false
			for nextNode := 0; nextNode < n; nextNode++ {
				if _, ok := visited[nextNode]; !ok && matrix[lastNode][nextNode] > 0 {
					newPath := Path{
						nodes:  append([]int{}, currentPath.nodes...),
						weight: currentPath.weight + matrix[lastNode][nextNode],
					}
					newPath.nodes = append(newPath.nodes, nextNode)
					queue = append(queue, newPath)
					needsToBeQueued = true
				}
			}

			if !needsToBeQueued {
				if currentPath.weight > maxWeight {
					maxWeight = currentPath.weight
				}
			}
			for k := range visited {
				delete(visited, k)
			}
		}
	}

	for start := 0; start < n; start++ {
		bfs(start)
	}
	return maxWeight
}

func calcUserGrade(matrix [][]int, userAnswer []int) int {
	sumWeight := 0
	lastIdx := userAnswer[0]

	for i := 1; i < len(userAnswer); i++ {
		sumWeight += matrix[lastIdx][userAnswer[i]]
		lastIdx = userAnswer[i]
	}
	return sumWeight
}

func validation(matrix [][]int, userAnswer []int) error {
	lenM := len(matrix)
	lenU := len(userAnswer)
	uniqueAnswers := make(map[int]struct{})

	if lenM == 0 {
		return errors.New("Граф пустой")
	}

	if lenU == 0 {
		return errors.New("Пользователь ничего не ответил")
	}

	if lenU > lenM {
		return errors.New("Некорректная длина ответа пользователя")
	}

	for _, node := range userAnswer {
		if (node > len(matrix)-1) || (node < 0) {
			return errors.New("Пользователь выбрал ответы, которые не были предложены")
		}
	}

	for _, node := range userAnswer {
		uniqueAnswers[node] = struct{}{}
	}
	if lenU != len(uniqueAnswers) {
		return errors.New("Номера ответов не уникальны")
	}

	for _, node := range matrix {
		if len(node) != lenM {
			return errors.New("Некорректный граф")
		}
	}

	for i := 0; i < lenM; i++ {
		for j := 0; j < lenM; j++ {
			if matrix[i][j] < 0 {
				return errors.New("Вес ребра в графе отрицательный")
			}
			if i == j && matrix[i][j] != 0 {
				return errors.New("Граф имеет петли")
			}
		}
	}
	return nil
}
