package service

import (
	"log"
	"shorters/domain"
	"shorters/repository"
	"shorters/service/cache"
	"shorters/service/dto"
	"time"

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

func (s *linkService) Store(url string, user domain.User) (*dto.StoreResponse, error) {
	var link domain.Link
	link.Key = shortid.MustGenerate()
	link.URL = url
	link.Creator = user.Email
	if user.Email == "" {
		link.CreatedTime = time.Now().Unix()
		link.ExpiredTime = time.Now().Add(time.Minute).Unix()
	}
	err := s.linkRepository.Store(&link)
	if err != nil {
		return nil, err
	}
	if user.Email != "" {
		s.linkCache.Put(link.Key, link.URL)
	}
	return &dto.StoreResponse{Key: link.Key}, nil
}
