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

// New returns an HTTP handler for creating a new student.
//
// It expects a JSON body containing "name", "email", and "age".
// Validates input using go-playground/validator,
// inserts the student into storage, and returns the generated ID.
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var student types.Student

		// Decode the JSON request body
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("empty request body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON: %w", err)))
			return
		}

		// Validate input fields
		validate := validator.New()
		if err := validate.Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// Create new student
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Student created successfully", slog.String("id", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]any{
			"status":  "success",
			"message": "student created successfully",
			"data":    lastId,
		})
	}
}

//
// ──────────────────────────────── GET STUDENT BY ID ────────────────────────────────
//

// GetById returns an HTTP handler that fetches a student by their ID.
//
// The URL must include the {id} path parameter, e.g. GET /api/students/1.
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

// GetList returns an HTTP handler that retrieves all students.
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Fetching all students")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "students fetched successfully",
			"data":    students,
		})
	}
}

//
// ──────────────────────────────── UPDATE STUDENT (PATCH) ────────────────────────────────
//

// UpdateById returns an HTTP handler that updates one or more fields of a student.
//
// Accepts a partial JSON body (PATCH). Only allowed fields ("name", "email", "age")
// are included in the update map. Unallowed fields are ignored.
// Example: PATCH /api/students/1
func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Updating student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student ID")))
			return
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		// Only allow certain fields to be updated
		allowed := map[string]bool{
			"name":  true,
			"email": true,
			"age":   true,
		}

		updates := map[string]any{}
		for k, v := range body {
			if allowed[k] {
				updates[k] = v
			}
		}

		student, err := storage.Update(intId, updates)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "student updated successfully",
			"data":    student,
		})
	}
}

//
// ──────────────────────────────── DELETE STUDENT BY ID ────────────────────────────────
//

// DeleteById returns an HTTP handler that deletes a student by their ID.
//
// The URL must include the {id} path parameter, e.g. DELETE /api/students/1.
// Returns how many rows were deleted (0 or 1).
func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Deleting student by ID", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid student ID")))
			return
		}

		rowsDeleted, err := storage.Delete(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]any{
			"status":  "success",
			"message": "student deleted successfully",
			"data":    rowsDeleted,
		})
	}
}
