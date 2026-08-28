package dto

type Slug struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}
type Slugs struct {
	Slugs []Slug `json:"subscriptions"`
}
