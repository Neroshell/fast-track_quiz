package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Question represents a quiz question with text, options, and the correct answer
type Question struct {
	ID      int      `json:"id"`      // Unique ID of the question
	Text    string   `json:"text"`    // Question text
	Options []string `json:"options"` // Multiple choice options for the question
	Correct int      `json:"correct"` // Index of the correct answer in the options array
}

// Submission represents the answers submitted by the user
type Submission struct {
	Answers []int `json:"answers"` // List of selected answers, where index matches question ID
}

// Predefined list of questions for the quiz
var questions = []Question{
	{ID: 1, Text: "What is the capital of France?", Options: []string{"Berlin", "Madrid", "Paris", "Rome"}, Correct: 2},
	{ID: 2, Text: "What is 2 + 2?", Options: []string{"3", "4", "5", "6"}, Correct: 1},
	{ID: 3, Text: "What is the color of the sky?", Options: []string{"Green", "Blue", "Red", "Yellow"}, Correct: 1},
	{ID: 4, Text: "Which planet is closest to the Sun?", Options: []string{"Earth", "Mars", "Mercury", "Venus"}, Correct: 2},
	{ID: 5, Text: "What is the boiling point of water?", Options: []string{"50°C", "100°C", "75°C", "200°C"}, Correct: 1},
}

// Variables to keep track of total quizzes taken and total correct answers
var (
	totalQuizzesTaken = 0 // Total number of quizzes completed by users
	totalCorrect      = 0 // Total number of correct answers across all quizzes
	mu                sync.Mutex // Mutex to safely update shared variables
)

// getQuestionsHandler handles HTTP requests for retrieving the list of questions
func getQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Send the list of questions as a JSON response
	json.NewEncoder(w).Encode(questions)
}

// submitAnswersHandler processes the user's quiz answers and calculates their score
func submitAnswersHandler(w http.ResponseWriter, r *http.Request) {
	var submission Submission

	// Decode the incoming JSON request into a Submission struct
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		// If decoding fails, return a Bad Request response
		http.Error(w, "Invalid submission", http.StatusBadRequest)
		return
	}

	// Calculate the number of correct answers
	correctCount := 0
	for i, answer := range submission.Answers {
		// Check if the submitted answer matches the correct one
		if questions[i].Correct == answer {
			correctCount++
		}
	}

	// Lock the mutex to safely update global stats
	mu.Lock()
	totalQuizzesTaken++          // Increment the total number of quizzes taken
	totalCorrect += correctCount // Increment the total number of correct answers
	mu.Unlock()                  // Unlock the mutex after updating

	// Calculate the comparison percentage of correct answers across all users
	comparison := (float64(totalCorrect) / float64(totalQuizzesTaken*len(questions))) * 100

	// Prepare the result data to send back to the user
	result := map[string]interface{}{
		"correct":    correctCount,        // Number of correct answers for this user
		"total":      len(questions),      // Total number of questions
		"comparison": fmt.Sprintf("%.2f", comparison), // Comparison percentage
	}
	// Send the result as a JSON response
	json.NewEncoder(w).Encode(result)
}

func main() {
	// Route to get the list of questions
	http.HandleFunc("/questions", getQuestionsHandler)

	// Route to submit answers and calculate the score
	http.HandleFunc("/submit", submitAnswersHandler)

	// Start the server and listen on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
