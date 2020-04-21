package main

import (
	"github.com/AcroManiac/micropic/internal/adapters/logger"
	"github.com/AcroManiac/micropic/internal/domain/entities"
	"github.com/AcroManiac/micropic/internal/domain/interfaces"
)

type CacheSender struct {
	client interfaces.CacheClient
}

// NewCacheSender constructor
func NewCacheSender(client interfaces.CacheClient) *CacheSender {
	return &CacheSender{client: client}
}

// Send previewer response to cache
func (s *CacheSender) Send(response *entities.Response, obj interface{}) {
	err := s.client.Save(response.Preview)
	if err != nil {
		logger.Error("error sending preview to cache", "error", err)
	}
}
