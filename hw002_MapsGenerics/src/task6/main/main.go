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
	objectGradeResult := make(map[string]map[int][]int)

	for _, result := range exam.Results {
		student := findStudentById(exam.Students, result.Student_Id)
		object := findObjectById(exam.Objects, result.Object_id)
		if _, ok := objectGradeResult[object.Name]; !ok {
			objectGradeResult[object.Name] = make(map[int][]int)
		}
		if _, ok := objectGradeResult[object.Name][student.Grade]; !ok {
			objectGradeResult[object.Name][student.Grade] = make([]int, 0, len(exam.Results)/len(exam.Objects)/2)
		}
		objectGradeResult[object.Name][student.Grade] = append(objectGradeResult[object.Name][student.Grade], result.Res)
	}
	displayTable(objectGradeResult)
}
func displayTable(m map[string]map[int][]int) {
	for objName, grades := range m {
		var sumResults, sumStudents int
		countRepeat := 17
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
		fmt.Printf(" %-9s| %-5s\n", objName, "Mean")
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))

		for grade, gradeArr := range grades {
			sum := Reduce(gradeArr, 0, func(a, b int) int {
				return a + b
			})
			mean := float64(sum) / float64(len(gradeArr))
			fmt.Printf("%3d grade | %.1f\n", grade, mean)
			sumResults += sum
			sumStudents += len(gradeArr)
		}
		meanAll := float64(sumResults) / float64(sumStudents)
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
		fmt.Printf(" %-9s| %-5.1f\n", "mean", meanAll)
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
	}
}
func Reduce[T any](s []T, init T, f func(a, b T) T) T {
	r := init
	for _, t := range s {
		r = f(t, r)
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
