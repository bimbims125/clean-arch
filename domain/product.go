package domain

type Product struct {
	ID          int
	Name        string
	Description string `json:"description,omitempty"`
	Price       float64
	ImageURL    string
	Stock       int
	Sold        int
	Category    Category
}
