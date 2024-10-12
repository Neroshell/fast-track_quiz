package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"github.com/spf13/cobra"
)

// Question struct represents a quiz question with options and the correct answer index
type Question struct {
	ID      int      `json:"id"`      // Unique identifier for the question
	Text    string   `json:"text"`    // The question text
	Options []string `json:"options"` // Possible answer choices
	Correct int      `json:"correct"` // The correct answer's index (zero-based)
}

// AnswerSubmission struct represents the user's submitted answers
type AnswerSubmission struct {
	Answers []int `json:"answers"` // Array of the user's selected answer indices (zero-based)
}

// getQuestions is the command handler that fetches quiz questions from the server
func getQuestions(cmd *cobra.Command, args []string) {
	// Make a GET request to retrieve the questions
	resp, err := http.Get("http://localhost:8080/questions")
	if err != nil {
		fmt.Println("Error fetching questions:", err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Decode the JSON response into a slice of Question structs
	var questions []Question
	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		fmt.Println("Error decoding questions:", err)
		return
	}

	// Display questions and collect user answers
	fmt.Println("Quiz Questions:")
	userAnswers := make([]int, len(questions)) // Initialize user answers array
	for i, q := range questions {
		fmt.Printf("%d. %s\n", q.ID, q.Text) // Display the question text
		for idx, option := range q.Options {
			fmt.Printf("  %d: %s\n", idx+1, option) // Display the available options
		}

		// Prompt the user to select an answer
		var answer int
		fmt.Printf("Enter your answer (1-%d): ", len(q.Options))
		_, err := fmt.Scan(&answer) // Read user's input
		if err != nil || answer < 1 || answer > len(q.Options) {
			fmt.Println("Invalid input. Please try again.")
			i-- // Ask the same question again if the input is invalid
			continue
		}
		userAnswers[i] = answer - 1 // Store user's answer (adjust for zero-based index)
	}

	// Submit user's answers to the server
	submitAnswers(userAnswers)
}

// submitAnswers sends the user's answers to the server and displays the result
func submitAnswers(answers []int) {
	// Create the submission object with user's answers
	submission := AnswerSubmission{Answers: answers}
	submissionJSON, err := json.Marshal(submission) // Convert to JSON format
	if err != nil {
		fmt.Println("Error creating submission:", err)
		return
	}

	// Send the submission as a POST request to the server
	resp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(submissionJSON))
	if err != nil {
		fmt.Println("Error submitting answers:", err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Decode the server's JSON response into a result map
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding result:", err)
		return
	}

	// Display the result to the user
	fmt.Printf("You got %v out of %v questions correct!\n", result["correct"], result["total"])
	fmt.Printf("You were better than %v%% of all quizzers.\n", result["comparison"])
}

func main() {
	// Define the root command for the CLI tool
	var rootCmd = &cobra.Command{Use: "quiz-cli"}

	// Define the get-questions command which fetches quiz questions
	var getCmd = &cobra.Command{
		Use:   "get-questions",
		Short: "Fetch quiz questions",
		Run:   getQuestions, // Attach the getQuestions handler
	}

	// Add the get-questions command to the root command
	rootCmd.AddCommand(getCmd)

	// Execute the root command (this runs the CLI application)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1) // Exit with status code 1 if there is an error
	}
}
