package database

import "github.com/samualhalder/go-restapis/internals/types"

type Database interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
