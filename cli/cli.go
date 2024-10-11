package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"github.com/spf13/cobra"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Correct int      `json:"correct"`
}

type AnswerSubmission struct {
	Answers []int `json:"answers"`
}

func getQuestions(cmd *cobra.Command, args []string) {
	resp, err := http.Get("http://localhost:8080/questions")
	if err != nil {
		fmt.Println("Error fetching questions:", err)
		return
	}
	defer resp.Body.Close()

	var questions []Question
	err = json.NewDecoder(resp.Body).Decode(&questions)
	if err != nil {
		fmt.Println("Error decoding questions:", err)
		return
	}

	// Display questions and options to the user
	fmt.Println("Quiz Questions:")
	userAnswers := make([]int, len(questions))
	for i, q := range questions {
		fmt.Printf("%d. %s\n", q.ID, q.Text)
		for idx, option := range q.Options {
			fmt.Printf("  %d: %s\n", idx+1, option)
		}

		// Capture user's answer
		var answer int
		fmt.Printf("Enter your answer (1-%d): ", len(q.Options))
		_, err := fmt.Scan(&answer)
		if err != nil || answer < 1 || answer > len(q.Options) {
			fmt.Println("Invalid input. Please try again.")
			i-- // Ask the same question again if input is invalid
			continue
		}
		userAnswers[i] = answer - 1 // Store answer (adjusting for zero-based index)
	}

	// Submit answers
	submitAnswers(userAnswers)
}

func submitAnswers(answers []int) {
	// Prepare the submission
	submission := AnswerSubmission{Answers: answers}
	submissionJSON, err := json.Marshal(submission)
	if err != nil {
		fmt.Println("Error creating submission:", err)
		return
	}

	// Send the POST request
	resp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(submissionJSON))
	if err != nil {
		fmt.Println("Error submitting answers:", err)
		return
	}
	defer resp.Body.Close()

	// Decode and display the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding result:", err)
		return
	}

	// Display the result
	fmt.Printf("You got %v out of %v questions correct!\n", result["correct"], result["total"])
	fmt.Printf("You were better than %v%% of all quizzers.\n", result["comparison"])
}

func main() {
	var rootCmd = &cobra.Command{Use: "quiz-cli"}

	var getCmd = &cobra.Command{
		Use:   "get-questions",
		Short: "Fetch quiz questions",
		Run:   getQuestions,
	}

	rootCmd.AddCommand(getCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
