package storage

import "github.com/gourav224/student-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	Update(id int64, updates map[string]any) (types.Student, error)
	Delete(id int64) (int64, error)
}
