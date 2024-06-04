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
	ObjectId  int `json:"object_id"`
	StudentId int `json:"student_id"`
	Res       int `json:"result"`
}

func main() {
	exam, err := parseJson("src/dz3.json")
	if err != nil {
		return
	}
	if len(exam.Students) == 0 {
		fmt.Println("Файл не содержит данных")
		return
	}

	mapStudents := getMapStudents(exam.Students)
	mapObjects := getMapObjects(exam.Objects)

	fmt.Printf("%-44s\n", strings.Repeat("_", 44))
	fmt.Printf("%-14s|%-7s|%-11s|%-9s\n", " Student name", " Grade", " Object", " Result")
	fmt.Printf("%-44s\n", strings.Repeat("_", 44))

	for _, r := range exam.Results {
		student := findStudentById(mapStudents, r.StudentId)
		object := findObjectById(mapObjects, r.ObjectId)
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

func findStudentById(mapStudents map[int]Student, studentId int) Student {
	return mapStudents[studentId]
}

func findObjectById(mapObjects map[int]Object, objectId int) Object {
	return mapObjects[objectId]
}

func getMapStudents(students []Student) map[int]Student {
	mapStudents := make(map[int]Student, len(students))
	for _, student := range students {
		mapStudents[student.Id] = student
	}
	return mapStudents
}

func getMapObjects(objects []Object) map[int]Object {
	mapObjects := make(map[int]Object, len(objects))
	for _, object := range objects {
		mapObjects[object.Id] = object
	}
	return mapObjects
}
