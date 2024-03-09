package main

import (
	"database/sql"
	"fmt"
	"golang-transaction/entity"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "001213"
	dbname   = "enigmacamp"
)

var psqlinfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func main() {
	studentEnrollment := entity.StudentEnrollment{Id: 2, Student_Id: 8, Subject: "AMBP", Credit: 4}

	enrollSubject(studentEnrollment)
}

func enrollSubject(studentEnrollment entity.StudentEnrollment) {
	db := connectDb()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	insertStudentEnrollment(studentEnrollment, tx)

	takenCredit := getSumCreditOfStudent(studentEnrollment.Student_Id, tx)

	updateStudent(takenCredit, studentEnrollment.Student_Id, tx)

	err = tx.Commit()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Transaction Committed!")
	}
}

func insertStudentEnrollment(studentEnrollment entity.StudentEnrollment, tx *sql.Tx) {
	insertStudentEnrollment := "INSERT INTO tx_student_enrollment (id, student_id, subject, credit) VALUES ($1, $2, $3, $4);"

	_, err := tx.Exec(insertStudentEnrollment, studentEnrollment.Id, studentEnrollment.Student_Id, studentEnrollment.Subject, studentEnrollment.Credit)
	validate(err, "Insert", tx)
}

func getSumCreditOfStudent(id int, tx *sql.Tx) int {
	sumCredit := "SELECT SUM(credit) FROM tx_student_enrollment WHERE student_id = $1;"

	takenCredit := 0
	err := tx.QueryRow(sumCredit, id).Scan(&takenCredit)
	validate(err, "Select", tx)

	return takenCredit
}

func updateStudent(takenCredit int, studentId int, tx *sql.Tx) {
	updateStudent := "UPDATE mst_student SET taken_credit = $1 WHERE id = $2"

	_, err := tx.Exec(updateStudent, takenCredit, studentId)
	validate(err, "Update", tx)
}

func validate(err error, message string, tx *sql.Tx) {
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully" + message + "data!")
	}
}

func connectDb() *sql.DB {
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Succesfully Connected!")
	}

	return db
}
