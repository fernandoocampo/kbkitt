package kbs

import (
	"context"
	"fmt"
	"strings"
)

type KBServiceClient interface {
	Create(ctx context.Context, newKB NewKB) (string, error)
	Import(ctx context.Context, newKBs []NewKB) (*ImportResult, error)
	Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error)
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
		return nil, fmt.Errorf("failed to add kb: %w", err)
	}

	kb := newKB.toKB(id)
	return &kb, nil
}

func (s *Service) Import(ctx context.Context, newKBs []NewKB) (*ImportResult, error) {
	err := validateKBs(newKBs)
	if err != nil {
		return nil, fmt.Errorf("one kb is not valid: %w", err)
	}

	result, err := s.kbClient.Import(ctx, newKBs)
	if err != nil {
		return nil, fmt.Errorf("failed to add kb: %w", err)
	}

	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*KB, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("the given id is not valid, because it is empty")
	}

	kb, err := s.kbClient.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get kb: %w", err)
	}

	return kb, nil
}

func (s *Service) Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error) {
	if IsStringEmpty(filter.Key) && IsStringEmpty(filter.Keyword) {
		return nil, nil
	}

	kbs, err := s.kbClient.Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search kb: %w", err)
	}

	return kbs, nil
}

func validateKBs(newKBs []NewKB) error {
	for index, newKB := range newKBs {
		err := newKB.validate()
		if err != nil {
			return fmt.Errorf("record %d in file is not valid: %w", index, err)
		}
	}

	return nil
}
