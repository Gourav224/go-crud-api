package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/gourav224/student-api/internal/config"
	"github.com/gourav224/student-api/internal/types"
	_ "github.com/mattn/go-sqlite3" // Import for side effects: registers the SQLite3 driver
)

// Sqlite wraps the SQL database connection.
type Sqlite struct {
	Db *sql.DB
}

// New initializes and returns a new SQLite connection.
// It also ensures the 'students' table exists before returning.
func New(cfg *config.Config) (*Sqlite, error) {
	// Open database file (creates if not exists)
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Verify the database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	// Create the students table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);`

	if _, err = db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("failed to create students table: %w", err)
	}

	return &Sqlite{Db: db}, nil
}

// CreateStudent inserts a new student record into the 'students' table.
// Returns the ID of the newly created student.
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// Prepare the INSERT statement
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the statement with provided parameters
	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	// Retrieve the last inserted ID
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

// GetStudentById retrieves a single student record by its ID.
// Returns a Student struct or an error if not found.
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	// Prepare the SELECT statement
	stmt, err := s.Db.Prepare("SELECT id, email, name, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	// Query a single row and scan the result into the student struct
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Email, &student.Name, &student.Age)
	if err != nil {
		return types.Student{}, err
	}

	return student, nil
}

// GetStudents retrieves all student records from the 'students' table.
// Returns a slice of Student structs or an error.
func (s *Sqlite) GetStudents() ([]types.Student, error) {
	// Prepare the SELECT statement
	stmt, err := s.Db.Prepare("SELECT id, email, name, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query to get multiple rows
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	// Iterate over the result set and map each row to a Student struct
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Email, &student.Name, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	// Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// Update modifies one or more fields of a student record.
// Accepts a map[string]any so the user can update a single field or multiple fields.
// Builds a dynamic SQL UPDATE statement using only the provided fields.
// Returns the updated student or an error if the student does not exist or update fails.
func (s *Sqlite) Update(id int64, updates map[string]any) (types.Student, error) {

	// Ensure at least one field is being updated
	if len(updates) == 0 {
		return types.Student{}, fmt.Errorf("no fields to update")
	}

	// Check if student exists
	_, err := s.GetStudentById(id)
	if err != nil {
		return types.Student{}, err
	}

	// Start building dynamic UPDATE query
	query := "UPDATE students SET "

	args := []any{}
	i := 0

	// Add each field to SQL query
	for k, v := range updates {
		if i > 0 {
			query += ", "
		}
		query += k + " = ?"
		args = append(args, v)
		i++
	}

	// Add WHERE clause
	query += " WHERE id = ?"
	args = append(args, id)

	// Prepare the dynamic UPDATE statement
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	// Execute UPDATE with values
	_, err = stmt.Exec(args...)
	if err != nil {
		return types.Student{}, err
	}

	// Return updated student
	return s.GetStudentById(id)
}

// Delete removes a student by ID from the database.
// Returns the number of rows deleted (0 or 1).
func (s *Sqlite) Delete(id int64) (int64, error) {
	// Ensure the student exists before deleting
	_, err := s.GetStudentById(id)
	if err != nil {
		return 0, fmt.Errorf("student not found: %w", err)
	}

	// Prepare DELETE query
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the delete
	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	// How many rows were deleted?
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
