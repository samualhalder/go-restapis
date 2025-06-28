package database

type Database interface {
	CreateStudent(id int, name string, email string, age int) (int, error)
}
