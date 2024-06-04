package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Cache[K comparable, V any] struct {
	m map[K]V
}

func (c *Cache[K, V]) Init() {
	c.m = make(map[K]V)
}
func (c *Cache[K, V]) Set(key K, value V) {
	c.m[key] = value
}
func (c *Cache[K, V]) Get(key K) (V, bool) {
	k, ok := c.m[key]
	return k, ok
}

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
	ObjectId  int `json:"object_id"`
	StudentId int `json:"student_id"`
	Res       int `json:"result"`
}

var students Cache[int, *Student]
var objects Cache[int, *Object]
var results Cache[int, *Result]

func main() {
	exam, err := parseJson("src/dz3.json")
	if err != nil {
		return
	}
	if len(exam.Students) == 0 {
		fmt.Println("Файл не содержит данных")
		return
	}

	students.Init()
	objects.Init()
	results.Init()

	for _, student := range exam.Students {
		students.Set(student.Id, &student)
	}
	for _, object := range exam.Objects {
		objects.Set(object.Id, &object)
	}

	fmt.Printf("%-44s\n", strings.Repeat("_", 44))
	fmt.Printf("%-14s|%-7s|%-11s|%-9s\n", " Student name", " Grade", " Object", " Result")
	fmt.Printf("%-44s\n", strings.Repeat("_", 44))

	for _, r := range exam.Results {
		student, ok := students.Get(r.StudentId)
		if !ok {
			defer fmt.Printf("Студента с Id = %d нет в списке студентов\n", r.StudentId)
			continue
		}
		object, ok := objects.Get(r.ObjectId)
		if !ok {
			defer fmt.Printf("Предмета с Id = %d нет в списке предметов\n", r.ObjectId)
			continue
		}

		studentName := student.Name
		studentGrade := student.Grade
		objectName := object.Name

		fmt.Printf(" %-13s|%5d  | %-10s|  %-7d\n", studentName, studentGrade, objectName, r.Res)
	}
}

func parseJson(filePath string) (Exam, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла", err)
		return Exam{}, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Ошибка закрытия файла")
		}
	}(file)

	decoder := json.NewDecoder(file)
	var exam Exam
	err2 := decoder.Decode(&exam)
	if err2 != nil {
		fmt.Println("Ошибка декодирования файла", err2)
		return Exam{}, err2
	}
	return exam, nil
}
