package service

import (
	"shorters/domain"
	"shorters/service/dto"
)

type LinkService interface {
	Find(key string) (string, error)
	Shorten(url string, user domain.User) (*dto.ShortenResponse, error)
	CustomShorten(url, customKey string, user domain.User) (*dto.ShortenResponse, error)
}
