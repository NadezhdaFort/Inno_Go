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
	votesArr := make(map[string]int)
	for _, name := range candidates {
		votesArr[name]++
	}
	var result []Candidate
	for name, votes := range votesArr {
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
