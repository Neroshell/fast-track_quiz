package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Correct int      `json:"correct"`
}

type Submission struct {
	Answers []int `json:"answers"`
}

var questions = []Question{
	{ID: 1, Text: "What is the capital of France?", Options: []string{"Berlin", "Madrid", "Paris", "Rome"}, Correct: 2},
	{ID: 2, Text: "What is 2 + 2?", Options: []string{"3", "4", "5", "6"}, Correct: 1},
	{ID: 3, Text: "What is the color of the sky?", Options: []string{"Green", "Blue", "Red", "Yellow"}, Correct: 1},
	{ID: 4, Text: "Which planet is closest to the Sun?", Options: []string{"Earth", "Mars", "Mercury", "Venus"}, Correct: 2},
	{ID: 5, Text: "What is the boiling point of water?", Options: []string{"50°C", "100°C", "75°C", "200°C"}, Correct: 1},
}

var (
	totalQuizzesTaken = 0
	totalCorrect      = 0
	mu                sync.Mutex
)

func getQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(questions)
}

func submitAnswersHandler(w http.ResponseWriter, r *http.Request) {
	var submission Submission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, "Invalid submission", http.StatusBadRequest)
		return
	}

	// Calculate score
	correctCount := 0
	for i, answer := range submission.Answers {
		if questions[i].Correct == answer {
			correctCount++
		}
	}

	// Update total quizzes taken and total correct answers
	mu.Lock()
	totalQuizzesTaken++
	totalCorrect += correctCount
	mu.Unlock()

	// Calculate comparison percentage
	comparison := (float64(totalCorrect) / float64(totalQuizzesTaken*len(questions))) * 100

	// Send response
	result := map[string]interface{}{
		"correct":    correctCount,
		"total":      len(questions),
		"comparison": fmt.Sprintf("%.2f", comparison),
	}
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/questions", getQuestionsHandler)
	http.HandleFunc("/submit", submitAnswersHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
