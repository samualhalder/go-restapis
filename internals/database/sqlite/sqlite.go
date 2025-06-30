package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samualhalder/go-restapis/internals/config"
	"github.com/samualhalder/go-restapis/internals/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no Such user for id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("error is %s", err.Error())
	}
	return student, nil
}

func (s *Sqlite) GetStudentList() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT id,name,email,age FROM students")
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return []types.Student{}, err
		}
		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateSutdent(id int64, name string, email string, age int) (bool, error) {
	stmt, err := s.Db.Prepare("UPDATE students SET name=?,email=?,age=? WHERE id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no student found for id %s", fmt.Sprint(id))
		}
		return false, fmt.Errorf("error %s", err.Error())
	}
	_, err = result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error %s", err.Error())
	}

	return true, nil
}
func (s *Sqlite) DeleteSutdent(id int64) (bool, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no student found for id %s", fmt.Sprint(id))
		}
		return false, fmt.Errorf("error %s", err.Error())
	}
	_, err = result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error %s", err.Error())
	}

	return true, nil
}
