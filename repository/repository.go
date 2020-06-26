package repository

import "shorters/domain"

type LinkRepository interface {
	Find(key string) (*domain.Link, error)
	FindByUser(email string) ([]*domain.Link, error)
	Store(link *domain.Link) error
	AddVisits(key string) error
}

type UserRepository interface {
	FindOTP(email string) (string, error)
}
