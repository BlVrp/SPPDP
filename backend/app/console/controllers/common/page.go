package common

// Page pagination data container response.
type Page[T any] struct {
	Data  []T `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
