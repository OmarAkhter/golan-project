package sqlite

import (
	"database/sql"

	"github.com/OmarAkhter/golan-project/internal/config"
	"github.com/OmarAkhter/golan-project/internal/types"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	// Add necessary fields like DB connection here
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL
	);
	`)
	if err != nil {
		return nil, err
	}

	return &SQLite{Db: db}, nil

}

func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {
	// Implementation for creating a student in SQLite
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SQLite) GetStudentByID(id int64) (types.Student, error) {
	// Implementation for retrieving a student by ID from SQLite
	var student types.Student
	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		return types.Student{}, err
	}
	return student, nil
}

func (s *SQLite) GetStudentList() ([]types.Student, error) {
	// Implementation for retrieving a list of students from SQLite
	rows, err := s.Db.Query("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *SQLite) DeleteStudentByID(id int64) error {
	// Implementation for deleting a student by ID from SQLite
	_, err := s.Db.Exec("DELETE FROM students WHERE id = ?", id)
	return err
}
