package domain

import "fmt"

type LinkNotFoundError struct {
	Key string
}

func (e LinkNotFoundError) Error() string {
	return fmt.Sprintf("no link found with key: %v", e.Key)
}
