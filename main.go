package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
)

type Learner struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
    Summary string `json:"summary,omitempty"`
}

var (
    learnerRecords = make(map[int]Learner)
    recordLock     sync.Mutex
    currentID      = 1
)

// Register a new learner
func registerLearner(w http.ResponseWriter, r *http.Request) {
    var learner Learner
    json.NewDecoder(r.Body).Decode(&learner)

    recordLock.Lock()
    learner.ID = currentID
    learnerRecords[currentID] = learner
    currentID++
    recordLock.Unlock()

    json.NewEncoder(w).Encode(learner)
}

// Retrieve all learners
func listLearners(w http.ResponseWriter, r *http.Request) {
    recordLock.Lock()
    defer recordLock.Unlock()

    var learners []Learner
    for _, learner := range learnerRecords {
        learners = append(learners, learner)
    }
    json.NewEncoder(w).Encode(learners)
}

// Retrieve learner details by ID
func fetchLearnerByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/learners/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    recordLock.Lock()
    learner, found := learnerRecords[id]
    recordLock.Unlock()

    if !found {
        http.Error(w, "Learner not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(learner)
}

// Create or retrieve a summary for a learner by ID
func fetchLearnerSummary(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/learners/summary/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    recordLock.Lock()
    learner, found := learnerRecords[id]
    recordLock.Unlock()

    if !found {
        http.Error(w, "Learner not found", http.StatusNotFound)
        return
    }

    if learner.Summary == "" {
        learner.Summary = generateSummary(learner)
        recordLock.Lock()
        learnerRecords[id] = learner
        recordLock.Unlock()
    }

    json.NewEncoder(w).Encode(map[string]string{"summary": learner.Summary})
}

// Summary generation helper
func generateSummary(learner Learner) string {
    return fmt.Sprintf("Learner %s, Age: %d, Contact: %s", learner.Name, learner.Age, learner.Email)
}

// Update learner information by ID
func updateLearnerByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/learners/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var updatedLearner Learner
    json.NewDecoder(r.Body).Decode(&updatedLearner)

    recordLock.Lock()
    learner, found := learnerRecords[id]
    if found {
        updatedLearner.ID = learner.ID
        learnerRecords[id] = updatedLearner
    }
    recordLock.Unlock()

    if !found {
        http.Error(w, "Learner not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(updatedLearner)
}

// Delete a learner by ID
func deleteLearnerByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/learners/"):])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    recordLock.Lock()
    _, found := learnerRecords[id]
    if found {
        delete(learnerRecords, id)
    }
    recordLock.Unlock()

    if !found {
        http.Error(w, "Learner not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func main() {
    http.HandleFunc("/learners", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            registerLearner(w, r)
        case "GET":
            listLearners(w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/learners/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            fetchLearnerByID(w, r)
        case "PUT":
            updateLearnerByID(w, r)
        case "DELETE":
            deleteLearnerByID(w, r)
        default:
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/learners/summary/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            fetchLearnerSummary(w, r)
        } else {
            http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Server listening on port 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
