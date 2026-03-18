package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Date        string `json:"created_at"`
}

func get(url string) error {
	parts := strings.Split(url, "/")
	link := fmt.Sprintf("https://api.github.com/repos/%s/%s", parts[3], parts[4])

	resp, err := http.Get(link)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API return non-OK status: %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var repo Repo
	err = json.Unmarshal(body, &repo)
	if err != nil {
		log.Fatalf("Error retrieving from JSON: %v", err)
	}
	repo.Date = strings.ReplaceAll(repo.Date, "T", " ")
	repo.Date = strings.ReplaceAll(repo.Date, "Z", "")
	fmt.Printf("---------Information---------\n")
	fmt.Printf("Name of the repository: %v\n", repo.Name)
	fmt.Printf("Description: %v\n", repo.Description)
	fmt.Printf("Stars: %v\n", repo.Stars)
	fmt.Printf("Forks: %v\n", repo.Forks)
	fmt.Printf("Created: %v\n", repo.Date)
	fmt.Printf("-----------------------------\n")
	return nil
}

func main() {
	run := flag.NewFlagSet("get", flag.ExitOnError)
	if len(os.Args) < 3 {
		fmt.Printf("Expected appy {COMMAND} {ARG}\n")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "get":
		if err := get(os.Args[2]); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
		run.Parse(os.Args[3:])
	default:
		fmt.Printf("Unknown command: %v", os.Args[1])
		os.Exit(1)
	}
}
