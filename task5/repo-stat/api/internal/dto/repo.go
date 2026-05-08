package dto

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Date        string `json:"created_at"`
}

type Repos struct {
	Repos []Repo `json:"repositories"`
}
