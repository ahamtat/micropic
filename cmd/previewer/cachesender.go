package main

import "github.com/AcroManiac/micropic/internal/domain/entities"

type CacheSender struct {
	//
}

// NewCacheSender constructor
func NewCacheSender() *CacheSender {
	return &CacheSender{}
}

func (s *CacheSender) Send(response *entities.Response, obj interface{}) {
	//
}
