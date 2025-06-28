package students

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/samualhalder/go-restapis/internals/types"
	"github.com/samualhalder/go-restapis/internals/utils/response"
)

func New() http.HandlerFunc {
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
		response.ResponseWriter(w, http.StatusCreated, response.SuccessWriter(student, "Student created successfully"))
	}
}
