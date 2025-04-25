package models

type Note struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Added    string `json:"added"`
	Modified string `json:"modified"`
}
