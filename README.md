# Quiz App

This is a simple quiz application that includes a Command Line Interface (CLI) client and a REST API server, both developed in Go. The application allows users to take quizzes with multiple-choice questions, submit answers, and see how their results compare to others.

## Features

- Fetch questions with a set number of multiple-choice answers.
- Allow users to select one answer per question.
- Enable users to submit their answers and see how many they got correct.
- Display users' percentile ranking compared to other quiz takers.

## Requirements

- Go 1.16 or higher.

## Installation

To install and run this project, follow these steps:

1. **Clone the repository:**

   git clone https://github.com/your_username/quiz-app.git
   
2. **Navigate to the project directory:**

cd quiz-app

3. **Install dependencies:**

go mod tidy

## Usage

**Starting the Server**
To start the server, navigate to the server directory and run the following command:

cd server
go run main.go

The server will start on http://localhost:8080 and will be ready to accept requests.

**Running the CLI Client**
Open another terminal window, navigate to the cmd/cli directory, and execute:

cd cmd/cli
go run main.go

Follow the on-screen prompts to answer the quiz questions. The client will interact with the server to fetch questions and submit your answers.

## Additional Features

- Input Validation: Ensures that users input only valid answer numbers.

- Repeat Quiz: Allows users to retake the quiz without restarting the application.

- Detailed Statistics: Provides statistics such as the average time per question and the distribution of correct/incorrect answers.
