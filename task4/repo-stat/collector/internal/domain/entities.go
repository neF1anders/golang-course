package domain

import "time"

type Repo struct {
	Name        string
	Description string
	Stars       int
	Forks       int
	Date        time.Time
}

type Slug struct {
	Owner string
	Repo  string
}
