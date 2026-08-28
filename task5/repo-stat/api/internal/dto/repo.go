package dto

import "time"

type Repo struct {
	Name        string    `json:"full_name"`
	Description string    `json:"description"`
	Stars       int       `json:"stars"`
	Forks       int       `json:"forks"`
	Date        time.Time `json:"created_at"`
}

type Repos struct {
	Repos []Repo `json:"repositories"`
}
