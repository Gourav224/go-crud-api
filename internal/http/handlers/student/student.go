package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gourav224/student-api/internal/storage"
	"github.com/gourav224/student-api/internal/types"
	"github.com/gourav224/student-api/internal/utils/response"
)

//
// ──────────────────────────────── CREATE STUDENT ────────────────────────────────
//

// New returns an HTTP handler that creates a new student record.
//
// It expects a JSON body with "name", "email", and "age" fields.
// Example request body:
//
//	{
//	  "name": "John Doe",
//	  "email": "john@example.com",
//	  "age": 21
//	}
//
// Validates input using go-playground/validator, inserts into storage,
// and returns the ID of the newly created student.
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close() // Always close the request body

		var student types.Student

		// Decode JSON body into struct
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("empty request body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON: %w", err)))
			return
		}

		// Validate struct fields using go-playground/validator
		validate := validator.New()
		if err := validate.Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// Insert new record
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Student created successfully", slog.String("id", fmt.Sprint(lastId)))

		// Respond with success
		response.WriteJson(w, http.StatusCreated, map[string]any{
			"status":  "success",
			"message": "student created successfully",
			"data":    lastId,
		})
	}
}

//
// ──────────────────────────────── GET BY ID ────────────────────────────────
//

// GetById returns an HTTP handler that fetches a student by their ID.
//
// URL pattern should include an `{id}` path parameter.
// Example: GET /api/students/1
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Fetching student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student ID")))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Respond with student data
		response.WriteJson(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "student fetched successfully",
			"data":    student,
		})
	}
}

//
// ──────────────────────────────── GET ALL STUDENTS ────────────────────────────────
//

// GetList returns an HTTP handler that retrieves all students from the database.
//
// Example: GET /api/students
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Fetching all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Respond with list of students
		response.WriteJson(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "students fetched successfully",
			"data":    students,
		})
	}
}
