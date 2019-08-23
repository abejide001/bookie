package models

// Book type
type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
	Year   int    `json:year`
}
