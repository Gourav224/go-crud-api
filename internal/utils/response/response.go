package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: "error",
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is required", err.Field()))
		case "email":
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' must be a valid email", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' is invalid", err.Field()))
		}
	}

	return Response{
		Status: "error",
		Error:  strings.Join(errMsgs, ", "),
	}
}
