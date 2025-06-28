package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

func ResponseWriter(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	return err
}

func ErrorWriter(err error) Response {
	return Response{
		Error:   true,
		Message: err.Error(),
		Result:  nil,
	}
}

func ValidatorError(errs validator.ValidationErrors) Response {
	var errMessage []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessage = append(errMessage, fmt.Sprintf("%s field is required", err.Field()))
		default:
			errMessage = append(errMessage, fmt.Sprintf("%s field is invalid", err.Field()))
		}
	}
	return Response{
		Error:   true,
		Message: strings.Join(errMessage, ","),
		Result:  nil,
	}
}

func SuccessWriter(data any, message string) Response {
	return Response{
		Error:   false,
		Message: message,
		Result:  data,
	}
}
