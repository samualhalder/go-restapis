package students

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/samualhalder/go-restapis/internals/database"
	"github.com/samualhalder/go-restapis/internals/types"
	"github.com/samualhalder/go-restapis/internals/utils/response"
)

func New(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.ResponseWriter(w, http.StatusBadRequest, response.ErrorWriter(err))
			return
		}
		if err != nil {
			response.ResponseWriter(w, http.StatusBadRequest, response.ErrorWriter(err))
		}

		if err := validator.New().Struct(student); err != nil {
			validationError := err.(validator.ValidationErrors)
			response.ResponseWriter(w, http.StatusBadRequest, response.ValidatorError(validationError))
			return
		}
		lastId, err := db.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.ResponseWriter(w, http.StatusBadRequest, response.ErrorWriter(err))
		}

		response.ResponseWriter(w, http.StatusCreated, response.SuccessWriter(lastId, "Student created successfully"))
	}
}

func GetUserById(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.ResponseWriter(w, http.StatusBadRequest, response.ErrorWriter(err))
			return
		}
		student, err := db.GetStudentById(id64)
		if err != nil {
			response.ResponseWriter(w, http.StatusBadRequest, response.ErrorWriter(err))
			return
		}
		response.ResponseWriter(w, http.StatusOK, response.SuccessWriter(student, ""))
	}
}
