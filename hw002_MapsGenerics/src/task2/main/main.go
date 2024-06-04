package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

func countVotes(candidates []string) []Candidate {
	if len(candidates) == 0 {
		fmt.Print("Входной массив отсутствует ")
		return nil
	}
	votesMap := make(map[string]int)
	for _, name := range candidates {
		votesMap[name]++
	}
	result := make([]Candidate, 0, len(votesMap))
	for name, votes := range votesMap {
		result = append(result, Candidate{Name: name, Votes: votes})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Votes > result[j].Votes
	})
	return result
}
func main() {
	result := countVotes([]string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"})
	fmt.Println(result)
}
