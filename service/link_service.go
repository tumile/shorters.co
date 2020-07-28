package service

import (
	"fmt"
	"log"
	"shorters/domain"
	"shorters/repository"
	"shorters/service/cache"
	"shorters/service/dto"

	"github.com/teris-io/shortid"
)

type linkService struct {
	linkCache      cache.LFUCache
	linkRepository repository.LinkRepository
}

func NewLinkService(linkRepository repository.LinkRepository) LinkService {
	return &linkService{cache.NewLFUCache(100), linkRepository}
}

func (s *linkService) Find(key string) (string, error) {
	url := s.linkCache.Get(key)
	if url == nil {
		link, err := s.linkRepository.Find(key)
		if err != nil {
			return "", err
		}
		s.linkCache.Put(link.Key, link.URL)
		url = link.URL
	}
	go func() {
		err := s.linkRepository.AddVisits(key)
		if err != nil {
			log.Println(err)
		}
	}()
	return url.(string), nil
}

func (s *linkService) addVisits(key string) {
	err := s.linkRepository.AddVisits(key)
	if err != nil {
		log.Println(err)
	}
}

func (s *linkService) Shorten(url string, user domain.User) (*dto.ShortenResponse, error) {
	var link domain.Link
	link.Key = shortid.MustGenerate()
	link.URL = url
	err := s.linkRepository.Store(&link)
	if err != nil {
		return nil, err
	}
	s.linkCache.Put(link.Key, link.URL)
	return &dto.ShortenResponse{Key: link.Key}, nil
}

func (s *linkService) CustomShorten(url, customKey string, user domain.User) (*dto.ShortenResponse, error) {
	_, err := s.linkRepository.Find(customKey)
	if err == nil {
		return nil, fmt.Errorf("Key exists")
	}
	var link domain.Link
	link.Key = customKey
	link.URL = url
	err = s.linkRepository.Store(&link)
	if err != nil {
		return nil, err
	}
	s.linkCache.Put(link.Key, link.URL)
	return &dto.ShortenResponse{Key: link.Key}, nil
}
