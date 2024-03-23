package main

import (
	"testing"
)

func TestFetchTodoByID(t *testing.T) {
	// Test fetching a valid TODO
	validTodo, err := fetchTodoByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if validTodo.ID != 1 {
		t.Errorf("Expected TODO ID to be 1, got %d", validTodo.ID)
	}

	// Test fetching an invalid TODO
	_, err = fetchTodoByID(-1)
	if err == nil {
		t.Error("Expected an error for fetching invalid TODO, got nil")
	}
}
