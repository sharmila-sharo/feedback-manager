package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Feedback struct {
	ID           int    `json:"id"`
	EmployeeID   string `json:"employeeId"`
	FeedbackText string `json:"feedbackText"`
	Rating       int    `json:"rating"`
}

var db *sql.DB

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS feedback (
        id SERIAL PRIMARY KEY,
        employee_id TEXT NOT NULL,
        feedback_text TEXT NOT NULL,
        rating INTEGER NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/feedback/", updateFeedbackHandler)
	http.HandleFunc("/feedback/delete/", deleteFeedbackHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	fmt.Fprintf(w, "Hello, Feedback Manager!")
}

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, employee_id, feedback_text, rating FROM feedback")
		if err != nil {
			http.Error(w, "Failed to fetch", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var feedbacks []Feedback
		for rows.Next() {
			var fb Feedback
			if err := rows.Scan(&fb.ID, &fb.EmployeeID, &fb.FeedbackText, &fb.Rating); err == nil {
				feedbacks = append(feedbacks, fb)
			}
		}
		json.NewEncoder(w).Encode(feedbacks)

	case http.MethodPost:
		var newFeedback Feedback
		err := json.NewDecoder(r.Body).Decode(&newFeedback)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.QueryRow(
			"INSERT INTO feedback (employee_id, feedback_text, rating) VALUES ($1, $2, $3) RETURNING id",
			newFeedback.EmployeeID, newFeedback.FeedbackText, newFeedback.Rating,
		).Scan(&newFeedback.ID)

		if err != nil {
			http.Error(w, "Failed to insert", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(newFeedback)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated Feedback
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		"UPDATE feedback SET employee_id=$1, feedback_text=$2, rating=$3 WHERE id=$4",
		updated.EmployeeID, updated.FeedbackText, updated.Rating, id,
	)
	if err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	updated.ID = id
	json.NewEncoder(w).Encode(updated)
}

func deleteFeedbackHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM feedback WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Feedback with ID %d deleted", id)
}
