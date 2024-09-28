package kbs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/filesystems"
)

type KBServiceClient interface {
	Create(ctx context.Context, newKB NewKB) (string, error)
	Update(ctx context.Context, kb *KB) error
	Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error)
	Get(ctx context.Context, id string) (*KB, error)
}

type ServiceSetup struct {
	Name            string
	FileForSyncPath string
	KBClient        KBServiceClient
}

type Service struct {
	kbClient        KBServiceClient
	fileForSyncPath string
}

func NewService(settings ServiceSetup) *Service {
	newService := Service{
		kbClient:        settings.KBClient,
		fileForSyncPath: settings.FileForSyncPath,
	}

	return &newService
}

func (s *Service) Add(ctx context.Context, newKB NewKB) (*KB, error) {
	err := newKB.validate()
	if err != nil {
		return nil, NewDataError(fmt.Sprintf("the given values are not valid: %s", err))
	}

	id, err := s.kbClient.Create(ctx, newKB)
	if err != nil && errors.As(err, &ClientError{}) {
		return nil, NewDataError(fmt.Sprintf("unable to add kb due to given data: %s", err))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to add kb: %w", err)
	}

	kb := newKB.toKB(id)
	return &kb, nil
}

func (s *Service) Update(ctx context.Context, kb KB) error {
	err := kb.validate()
	if err != nil {
		return NewDataError(fmt.Sprintf("the given values are not valid: %s", err))
	}

	err = s.kbClient.Update(ctx, &kb)
	if err != nil && errors.As(err, &ClientError{}) {
		return NewDataError(fmt.Sprintf("unable to update kb due to given data: %s", err))
	}

	if err != nil {
		return fmt.Errorf("failed to update kb: %w", err)
	}

	return nil
}

func (s *Service) Import(ctx context.Context, newKBs []NewKB) (*ImportResult, error) {
	err := validateKBs(newKBs)
	if err != nil {
		return nil, fmt.Errorf("one kb is not valid: %w", err)
	}

	result := ImportResult{
		NewIDs:     make(map[string]string),
		FailedKeys: make(map[string]string),
	}

	for _, newKB := range newKBs {
		id, err := s.kbClient.Create(ctx, newKB)
		if err != nil {
			result.FailedKeys[newKB.Key] = err.Error()
			continue
		}

		result.NewIDs[newKB.Key] = id
	}

	return &result, nil
}

func (s *Service) SaveForSync(ctx context.Context, newKB NewKB) error {
	newKBYAML, err := newKB.toYAML()
	if err != nil {
		return fmt.Errorf("unable to save new kb for later sync: %w", err)
	}

	empty, err := filesystems.FileEmpty(s.fileForSyncPath)
	if err != nil {
		return fmt.Errorf("unable to check sync file: %w", err)
	}

	var content []byte

	if !empty {
		content = append(content, []byte("---\n")...)
	}

	content = append(content, newKBYAML...)

	err = filesystems.SaveOrAppendFile(s.fileForSyncPath, content)
	if err != nil {
		return fmt.Errorf("unable to save new kb for later sync: %w", err)
	}

	return nil
}

func (s *Service) Sync(ctx context.Context) (*SyncResult, error) {
	newKBs, err := loadSyncFile(s.fileForSyncPath)
	if err != nil {
		return nil, fmt.Errorf("unable to process synchronization: %w", err)
	}

	if len(newKBs) == 0 {
		return nil, nil
	}

	err = validateKBs(newKBs)
	if err != nil {
		return nil, fmt.Errorf("one kb is not valid: %w", err)
	}

	go func() {
		_ = filesystems.TruncateFile(s.fileForSyncPath)
	}()

	result := SyncResult{
		NewIDs:     make(map[string]string),
		FailedKeys: make(map[string]string),
	}

	for _, newKB := range newKBs {
		id, err := s.kbClient.Create(ctx, newKB)
		if err != nil {
			result.FailedKeys[newKB.Key] = err.Error()
			continue
		}

		result.NewIDs[newKB.Key] = id
	}

	return &result, nil
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
