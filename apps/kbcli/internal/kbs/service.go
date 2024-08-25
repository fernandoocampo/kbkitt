package kbs

import (
	"context"
	"fmt"
)

type KBServiceClient interface {
	Create(ctx context.Context, newKB NewKB) (string, error)
	Search(ctx context.Context, filter KBQueryFilter) ([]KBItem, error)
	Get(ctx context.Context, id string) (*KB, error)
}

type ServiceSetup struct {
	Name     string
	KBClient KBServiceClient
}

type Service struct {
	kbClient KBServiceClient
}

func NewService(settings ServiceSetup) *Service {
	newService := Service{
		kbClient: settings.KBClient,
	}

	return &newService
}

func (s *Service) Add(ctx context.Context, newKB NewKB) (*KB, error) {
	err := newKB.validate()
	if err != nil {
		return nil, fmt.Errorf("the given values are not valid: %w", err)
	}

	id, err := s.kbClient.Create(ctx, newKB)
	if err != nil {
		return nil, fmt.Errorf("unable to add kb: %w", err)
	}

	kb := newKB.toKB(id)
	return &kb, nil
}
