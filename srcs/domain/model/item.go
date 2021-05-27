package model

type ID interface{}

type Item struct {
	Id   ID     `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}
