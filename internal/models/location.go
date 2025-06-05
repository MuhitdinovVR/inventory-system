package models

type Location struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	Type    string `json:"type"`
}
