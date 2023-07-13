package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	// Handle flags
	timeLimit := flag.Int("time-limit", 30, "Specify the time limit for each question.")	

	flag.Parse()

	// open the csv file
	f, err := os.Open("problems.csv")
	if err != nil {
		log.Fatal(err)
	}

	// close the file at the end of the program
	defer f.Close()

	// read the csv values
	r := csv.NewReader(f)
	var answer string
	var score, questions int
	timeout := time.Duration(*timeLimit) * time.Second // Set timeout duration

	// Create a channel to receive input from the user
	inputCh := make(chan string)

	// Get the user to press enter before starting the timer
	fmt.Print("Press [ENTER] to start.")
	fmt.Scanln()

	go func() {
		for {
			rec, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			// process the record
			// show the question
			fmt.Print(rec[0], " = ")

			// Start timer for each question
			timer := time.AfterFunc(timeout, func() {
				fmt.Println("\nYou ran out of time lol")
				inputCh <- "" // Send empty string to input channel
			})

			// get the answer
			fmt.Scanln(&answer)

			// Stop the timer(input received)
			timer.Stop()

			if answer == rec[1] {
				score++
			}
			questions++	// Add to total number of questions

			// Send answer to channel
			inputCh <- answer
		}
		close(inputCh) // Close the channel
	}()

	// Listen for input or timeout
	for input := range inputCh {
		if input == "" {
			// No input received within timeout duration
			break
		}
	}

	fmt.Printf("You got %d out of %d correct!\n", score, questions)
}