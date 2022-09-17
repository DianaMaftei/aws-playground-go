package model

type Pet struct {
	Id        int       `json:"id"`
	Category  Category  `json:"category"`
	Name      string    `json:"name"`
	PhotoUrls []string  `json:"photoUrls"`
	Tags      []Tag     `json:"tags"`
	Status    PetStatus `json:"status"`
}

type PetStatus string

const (
	Available PetStatus = "available"
	Pending   PetStatus = "pending"
	Sold      PetStatus = "sold"
)
