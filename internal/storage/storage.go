package storage

import "github.com/OmarAkhter/golan-project/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudentList() ([]types.Student, error)
	DeleteStudentByID(id int64) error
}
