package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
)

type Student struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
    Summary string `json:"summary,omitempty"`
}

var (
    students = make(map[int]Student)
    mu       sync.Mutex
    nextID   = 1
)

// Create a new student
func createStudent(w http.ResponseWriter, r *http.Request) {
    var student Student
    json.NewDecoder(r.Body).Decode(&student)

    mu.Lock()
    student.ID = nextID
    students[nextID] = student
    nextID++
    mu.Unlock()

    json.NewEncoder(w).Encode(student)
}

// Get all students
func getStudents(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    var result []Student
    for _, student := range students {
        result = append(result, student)
    }
    json.NewEncoder(w).Encode(result)
}

// Get a student by ID
func getStudentByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/students/"):])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    mu.Lock()
    student, exists := students[id]
    mu.Unlock()

    if !exists {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(student)
}

// Generate a summary for a student by ID
func getStudentSummary(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/students/summary/"):])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    mu.Lock()
    student, exists := students[id]
    mu.Unlock()

    if !exists {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    // Generate summary if not already present
    if student.Summary == "" {
        student.Summary, _ = generateProfileSummary(student)
        mu.Lock()
        students[id] = student // Update student map with the summary
        mu.Unlock()
    }

    json.NewEncoder(w).Encode(map[string]string{"summary": student.Summary})
}

// Dummy function for summary generation
func generateProfileSummary(student Student) (string, error) {
    return fmt.Sprintf("Student %s, age %d, email %s", student.Name, student.Age, student.Email), nil
}

// Update a student by ID
func updateStudentByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/students/"):])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    var updatedStudent Student
    json.NewDecoder(r.Body).Decode(&updatedStudent)

    mu.Lock()
    student, exists := students[id]
    if exists {
        updatedStudent.ID = student.ID
        students[id] = updatedStudent
    }
    mu.Unlock()

    if !exists {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(updatedStudent)
}

// Delete a student by ID
func deleteStudentByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/students/"):])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    mu.Lock()
    _, exists := students[id]
    if exists {
        delete(students, id)
    }
    mu.Unlock()

    if !exists {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func main() {
    http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            createStudent(w, r)
        case "GET":
            getStudents(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            getStudentByID(w, r)
        case "PUT":
            updateStudentByID(w, r)
        case "DELETE":
            deleteStudentByID(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // Summary route
    http.HandleFunc("/students/summary/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            getStudentSummary(w, r)
        } else {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
