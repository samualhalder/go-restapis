package response

import (
	"encoding/json"
	"net/http"
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
