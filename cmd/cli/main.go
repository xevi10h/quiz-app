package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"

    "github.com/spf13/cobra"
)

type Question struct {
    ID       int      `json:"id"`
    Question string   `json:"question"`
    Options  []string `json:"options"`
}

type Submission struct {
    Answers map[int]int `json:"answers"`
    StartAt time.Time   `json:"start_at"`
}

type Result struct {
    CorrectAnswers int     `json:"correct_answers"`
    TotalQuestions int     `json:"total_questions"`
    Percentile     float64 `json:"percentile"`
    AverageTime    float64 `json:"average_time"`
}

var rootCmd = &cobra.Command{
    Use:   "quizcli",
    Short: "CLI to take a quiz",
    Run: func(cmd *cobra.Command, args []string) {
        for {
            takeQuiz()
            fmt.Print("\nDo you want to take the quiz again? (yes/no): ")
            var response string
            fmt.Scan(&response)
            if response != "yes" {
                break
            }
        }
    },
}

func takeQuiz() {
    resp, err := http.Get("http://localhost:8080/questions")
    if err != nil {
        fmt.Println("Error retrieving the questions:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var questions []Question
    json.Unmarshal(body, &questions)

    submission := Submission{
        Answers: make(map[int]int),
        StartAt: time.Now(),
    }

    for _, q := range questions {
        fmt.Printf("\nQuestion %d: %s\n", q.ID, q.Question)
        for i, option := range q.Options {
            fmt.Printf("  %d) %s\n", i+1, option)
        }

        var answer int
        for {
            fmt.Print("Your answer (number): ")
            fmt.Scan(&answer)
            if answer >= 1 && answer <= len(q.Options) {
                submission.Answers[q.ID] = answer - 1
                break
            }
            fmt.Println("Invalid choice, please enter a number listed above.")
        }
    }

    submitAnswers(submission)
}

func submitAnswers(submission Submission) {
    jsonData, _ := json.Marshal(submission)
    resp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error submitting the answers:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var result Result
    json.Unmarshal(body, &result)

    fmt.Printf("\nCorrect answers: %d out of %d\n", result.CorrectAnswers, result.TotalQuestions)
    fmt.Printf("You scored higher than %.2f%% of all participants.\n", result.Percentile)
    fmt.Printf("Your average time per question was %.2f seconds.\n", result.AverageTime)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
