package domain

type Repo struct {
	Name        string
	Description string
	Stars       int
	Forks       int
	Date        string
}

type Slug struct {
	Owner string
	Repo  string
}
