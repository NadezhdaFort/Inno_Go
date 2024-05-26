package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Exam struct {
	Students []Student `json:"students"`
	Objects  []Object  `json:"objects"`
	Results  []Result  `json:"results"`
}
type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}
type Object struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Result struct {
	Object_id  int `json:"object_id"`
	Student_Id int `json:"student_id"`
	Res        int `json:"result"`
}

func main() {
	exam := parseJson("src/dz3.json")
	fmt.Printf("%-44s\n", strings.Repeat("_", 44))
	fmt.Printf("%-14s|%-7s|%-11s|%-9s\n", " Student name", " Grade", " Object", " Result")
	fmt.Printf("%-44s\n", strings.Repeat("_", 44))
	ballsArr := make([]int, len(exam.Students)+1, len(exam.Students)+1)
	for _, r := range exam.Results {
		ballsArr[r.Student_Id] += r.Res
	}
	countObjects := len(exam.Objects)
	excellentStudents := Filter(ballsArr, func(i int) bool {
		if i == countObjects*5 {
			return true
		}
		return false
	})
	for _, indexIdSt := range excellentStudents {

		student := findStudentById(exam.Students, indexIdSt)
		//object := findObjectById(exam.Objects, r.Object_id)
		studentName := student.Name
		studentGrade := student.Grade
		for _, object := range exam.Objects {
			objectName := object.Name
			fmt.Printf(" %-13s|%5d  | %-10s|  %-7d\n", studentName, studentGrade, objectName, 5)
		}
	}
}

func Filter[T any](s []T, f func(T) bool) []int {
	var r []int
	for i, v := range s {
		if f(v) {
			r = append(r, i)
		}
	}
	return r
}

func parseJson(filePath string) Exam {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var exam Exam
	err2 := decoder.Decode(&exam)
	if err2 != nil {
		fmt.Println("Ошибка декодирования файла", err2)
	}
	return exam
}

func findStudentById(students []Student, studentId int) Student {
	var result Student
	for _, student := range students {
		if student.Id == studentId {
			result = student
		}
	}
	return result
}

func findObjectById(objects []Object, objectId int) Object {
	var result Object
	for _, object := range objects {
		if object.Id == objectId {
			result = object
		}

	}
	return result
}
