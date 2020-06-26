package service

import (
	"shorters/domain"
	"shorters/service/dto"
)

type LinkService interface {
	Find(key string) (string, error)
	Store(url string, user domain.User) (*dto.StoreResponse, error)
}
