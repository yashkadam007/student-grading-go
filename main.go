package main

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func convertRecordToStudent(record []string) (student, error) {
	var s student
	if len(record) != 7 {
		return student{}, errors.New("invalid record")
	}

	s.firstName = record[0]
	s.lastName = record[1]
	s.university = record[2]
	test1Score, err1 := strconv.Atoi(record[3])
	test2Score, err2 := strconv.Atoi(record[4])
	test3Score, err3 := strconv.Atoi(record[5])
	test4Score, err4 := strconv.Atoi(record[6])

	if err4 != nil || err3 != nil || err2 != nil || err1 != nil {
		return student{}, errors.New("invalid scores")
	}
	s.test1Score = test1Score
	s.test2Score = test2Score
	s.test3Score = test3Score
	s.test4Score = test4Score

	return s, nil
}

func parseCSV(filePath string) []student {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file"+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error in reading file"+filePath, err)
	}
	var students []student
	for _, v := range records[1:] {
		s, err := convertRecordToStudent(v)
		if err != nil {
			log.Fatal("Error in reading record", err)
		}
		students = append(students, s)
	}

	return students
}

func averageScore(s1, s2, s3, s4 int) float32 {
	avg := float32(s1+s2+s3+s4) / 4
	return avg
}

func calculateGrade(students []student) []studentStat {
	var gradedStudents []studentStat
	for _, v := range students {
		var studentStat studentStat
		var grade Grade
		studentStat.student = v
		finalScore := averageScore(v.test1Score, v.test2Score, v.test3Score, v.test4Score)
		switch {
		case finalScore < 35:
			grade = F
		case finalScore >= 35 && finalScore < 50:
			grade = C
		case finalScore >= 50 && finalScore < 70:
			grade = B
		case finalScore >= 70:
			grade = A
		}
		studentStat.finalScore = finalScore
		studentStat.grade = grade
		gradedStudents = append(gradedStudents, studentStat)
	}
	return gradedStudents
}

func findOverallTopper(gradedStudents []studentStat) studentStat {
	var topper studentStat
	var topperScore float32
	for _, v := range gradedStudents {
		if v.finalScore > topperScore {
			topper = v
			topperScore = v.finalScore
		}
	}
	return topper
}

func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	universityTopper := make(map[string]studentStat)
	for _, v := range gs {
		university := v.university
		val, ok := universityTopper[university]
		if !ok {
			universityTopper[university] = v
			continue
		}
		if v.finalScore > val.finalScore {
			universityTopper[university] = v
		}
	}
	return universityTopper
}
