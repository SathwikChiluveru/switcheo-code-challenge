package types

import (
	"fmt"
)

type Item struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Details string `json:"details"`
}

func NewItem(id, name, details string) Item {
	return Item{
		ID:      id,
		Name:    name,
		Details: details,
	}
}

func (i Item) String() string {
	return fmt.Sprintf("ID: %s\nName: %s\nDetails: %s", i.ID, i.Name, i.Details)
}
