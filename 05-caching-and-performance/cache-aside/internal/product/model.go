package product

import "time"

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Source string

const (
	SourceCache    Source = "CACHE"
	SourceDatabase Source = "DATABASE"
)

type Result struct {
	Product Product `json:"product"`
	Source  Source  `json:"source"`
}
