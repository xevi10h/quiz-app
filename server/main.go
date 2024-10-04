package main

import (
    "encoding/json"
    "log"
    "math"
    "net/http"
    "sync"
    "time"
)

type Question struct {
    ID       int      `json:"id"`
    Question string   `json:"question"`
    Options  []string `json:"options"`
}

type Submission struct {
    Answers  map[int]int `json:"answers"`
    StartAt  time.Time   `json:"start_at"`
}

type Result struct {
    CorrectAnswers int     `json:"correct_answers"`
    TotalQuestions int     `json:"total_questions"`
    Percentile     float64 `json:"percentile"`
    AverageTime    float64 `json:"average_time"`
}

var (
    questions = []Question{
        {ID: 1, Question: "What is the capital of France?", Options: []string{"Paris", "London", "Rome", "Berlin"}},
        {ID: 2, Question: "Which planet is the largest in the solar system?", Options: []string{"Mars", "Jupiter", "Saturn", "Venus"}},
        {ID: 3, Question: "Who painted the Mona Lisa?", Options: []string{"Leonardo da Vinci", "Pablo Picasso", "Vincent van Gogh", "Rembrandt"}},
        {ID: 4, Question: "What is the chemical symbol for water?", Options: []string{"O2", "H2O", "CO2", "NaCl"}},
        {ID: 5, Question: "Who wrote 'Macbeth'?", Options: []string{"William Shakespeare", "Charles Dickens", "Jane Austen", "Mark Twain"}},
        {ID: 6, Question: "What is the hardest natural substance on Earth?", Options: []string{"Gold", "Iron", "Diamond", "Quartz"}},
        {ID: 7, Question: "What is the main ingredient in guacamole?", Options: []string{"Tomato", "Avocado", "Onion", "Pepper"}},
        {ID: 8, Question: "What is the speed of light?", Options: []string{"299,792 km/s", "150,000 km/s", "1,000 km/s", "3,000 km/s"}},
        {ID: 9, Question: "What is the largest ocean on Earth?", Options: []string{"Atlantic Ocean", "Indian Ocean", "Arctic Ocean", "Pacific Ocean"}},
        {ID: 10, Question: "How many continents are there?", Options: []string{"Six", "Seven", "Five", "Eight"}},
    }

    correctAnswers = map[int]int{
        1: 0,  // Paris
        2: 1,  // Jupiter
        3: 0,  // Leonardo da Vinci
        4: 1,  // H2O
        5: 0,  // William Shakespeare
        6: 2,  // Diamond
        7: 1,  // Avocado
        8: 0,  // 299,792 km/s
        9: 3,  // Pacific Ocean
        10: 1, // Seven
    }

    mu          sync.Mutex
    totalUsers  int
    userResults []int
	totalResponseTime float64
)

func main() {
    http.HandleFunc("/questions", getQuestions)
    http.HandleFunc("/submit", submitAnswers)

    log.Println("Server listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(questions)
}

func submitAnswers(w http.ResponseWriter, r *http.Request) {
    var submission Submission
    if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    elapsedTime := time.Since(submission.StartAt).Seconds()
    correct := 0
    for qID, answerIndex := range submission.Answers {
        if correctAnswers[qID] == answerIndex {
            correct++
        }
    }

    mu.Lock()
    totalUsers++
    userResults = append(userResults, correct)
    totalResponseTime += elapsedTime
    averageTime := totalResponseTime / float64(totalUsers)
    percentile := calculatePercentile(correct)
    mu.Unlock()

    result := Result{
        CorrectAnswers: correct,
        TotalQuestions: len(questions),
        Percentile:     percentile,
        AverageTime:    averageTime,
    }

    json.NewEncoder(w).Encode(result)
}

func calculatePercentile(score int) float64 {
    count := 0
    for _, s := range userResults {
        if s <= score {
            count++
        }
    }
    percentile := (float64(count) / float64(len(userResults))) * 100
    return math.Round(percentile*100) / 100 // Round to two decimal places
}