Fast Track Quiz Application
This repository contains a simple quiz application built with Go, designed to allow users to answer a series of quiz questions via a CLI that communicates with a REST API backend.

Features
Fetch quiz questions from a backend REST API.
Answer each question with multiple choice options.
Submit answers and get feedback on the number of correct answers.
Compare your score with other quiz takers.

Stack
Backend: Golang (Go)
CLI Framework: spf13/cobra
Database: In-memory (no persistent storage)

Functionality
The CLI fetches questions from the backend server.
The user is prompted to answer the questions one by one.
After answering all the questions, the answers are submitted to the backend.
The backend evaluates the answers and returns:
The number of correct answers.
A comparison showing how well the user performed relative to other quiz takers.

Getting Started
Prerequisites
Go must be installed on your machine. 


git clone this repo
cd fast-track-quiz
Install dependencies:

The only external dependency is the cobra CLI framework. You can install it by running:

go get -u github.com/spf13/cobra
Running the Application
Step 1: Run the Backend Server
The backend server exposes two main routes:

GET /questions: Fetches the list of quiz questions.
POST /submit: Submits the user's answers and returns the result.
To run the backend server:

go run main.go
This starts the server at http://localhost:8080.

Step 2: Use the CLI to Interact with the Quiz
The CLI allows the user to fetch questions, submit answers, and get results.

To run the CLI:

go run cli.go get-questions
You will be prompted to answer each question in the terminal.

After answering all the questions, your answers will be automatically submitted to the backend, and you will receive a score and comparison result.

