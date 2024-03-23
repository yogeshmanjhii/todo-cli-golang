package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

const baseURL = "https://jsonplaceholder.typicode.com/todos/"

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func fetchTodoByID(todoID int) (Todo, error) {
	url := fmt.Sprintf("%s%d", baseURL, todoID)
	resp, err := http.Get(url)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to fetch TODO %d: %v", todoID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Todo{}, fmt.Errorf("failed to fetch TODO %d: status code %d", todoID, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to read response body for TODO %d: %v", todoID, err)
	}

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		return Todo{}, fmt.Errorf("failed to unmarshal JSON for TODO %d: %v", todoID, err)
	}

	return todo, nil
}

func main() {
	flag.Parse()

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Channel to receive the fetched todos
	todoCh := make(chan Todo)

	// Start fetching the first 20 even numbered TODOs
	for i := 2; i <= 40; i += 2 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			todo, err := fetchTodoByID(id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
			todoCh <- todo
		}(i)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(todoCh)
	}()

	// Print the fetched todos
	for todo := range todoCh {
		fmt.Printf("Title: %s\nCompleted: %v\n\n", todo.Title, todo.Completed)
	}

	fmt.Println("Done fetching todos.")
}

