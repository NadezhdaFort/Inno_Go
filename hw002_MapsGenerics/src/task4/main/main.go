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
type ObjectStatistic struct {
	SumAllResultsObject    int
	SumAllStudentsInObject int
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

	objectGradeResult := make(map[int]map[int]ObjectStatistic)

	for _, result := range exam.Results {
		student := findStudentById(mapStudents, result.StudentId)
		object := findObjectById(mapObjects, result.ObjectId)

		if _, ok := objectGradeResult[object.Id]; !ok {
			objectGradeResult[object.Id] = make(map[int]ObjectStatistic)
		}

		if _, ok := objectGradeResult[object.Id][student.Grade]; !ok {
			objectGradeResult[object.Id][student.Grade] = ObjectStatistic{}
		}

		currentObjectStatistic := objectGradeResult[object.Id][student.Grade]
		currentObjectStatistic.SumAllResultsObject += result.Res
		currentObjectStatistic.SumAllStudentsInObject++
		objectGradeResult[object.Id][student.Grade] = currentObjectStatistic
	}

	displayTable(objectGradeResult, mapObjects)
}
func displayTable(m map[int]map[int]ObjectStatistic, mapObjects map[int]Object) {
	for objectId, grades := range m {
		var sumResults, sumStudents int
		objectName := findObjectById(mapObjects, objectId).Name

		countRepeat := 17
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
		fmt.Printf(" %-9s| %-5s\n", objectName, "Mean")
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
		for grade, gradeArr := range grades {
			mean := Mean(gradeArr.SumAllResultsObject, gradeArr.SumAllStudentsInObject)
			fmt.Printf("%3d grade | %.1f\n", grade, mean)
			sumResults += gradeArr.SumAllResultsObject
			sumStudents += gradeArr.SumAllStudentsInObject
		}
		meanAll := Mean(sumResults, sumStudents)
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
		fmt.Printf(" %-9s| %-5.1f\n", "mean", meanAll)
		fmt.Printf("%14s\n", strings.Repeat("_", countRepeat))
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

func Mean(a, b int) float64 {
	return float64(a) / float64(b)
}
