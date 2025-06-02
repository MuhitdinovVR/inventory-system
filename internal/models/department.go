package models

type Department struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Location string  `json:"location"`
	HeadID   *int    `json:"head_id,omitempty"`
	HeadName *string `json:"head_name,omitempty"`
}
