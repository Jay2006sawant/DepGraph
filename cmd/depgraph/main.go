package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Verify required environment variables
	required := []string{"GITHUB_TOKEN"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
	}
}

func main() {
	fmt.Println("DepGraph - Multi-Repo Dependency Visualizer")
	fmt.Println("Starting dependency analysis...")
	
	// TODO: Initialize components
	// - GitHub API client
	// - Database connection
	// - Repository scanner
	// - Dependency analyzer
	// - Web server for visualization
} 