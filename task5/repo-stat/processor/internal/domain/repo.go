package domain

import "time"

type Repo struct {
	Name        string
	Description string
	Stars       int32
	Forks       int32
	Date        time.Time
}
