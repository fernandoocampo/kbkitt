package kbs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/filesystems"
	"github.com/fernandoocampo/kbkitt/apps/kbcli/internal/adapters/webs"
)

type Storage interface {
	Create(ctx context.Context, newKB KB) (string, error)
	GetByID(ctx context.Context, id string) (*KB, error)
	GetByKey(ctx context.Context, key string) (*KB, error)
	Update(ctx context.Context, kb *KB) error
	Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error)
	GetAll(ctx context.Context, filter KBQueryFilter) (*GetAllResult, error)
}

type KBServiceClient interface {
	Create(ctx context.Context, newKB NewKB) (string, error)
	Update(ctx context.Context, kb *KB) error
	Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error)
	Get(ctx context.Context, id string) (*KB, error)
}

type ServiceSetup struct {
	KBStorage       Storage
	KBClient        KBServiceClient
	Name            string
	FileForSyncPath string
	DirForMediaPath string
}

type Service struct {
	kbClient        KBServiceClient
	storage         Storage
	fileForSyncPath string
	dirForMediaPath string
}

func NewService(settings ServiceSetup) *Service {
	newService := Service{
		kbClient:        settings.KBClient,
		storage:         settings.KBStorage,
		fileForSyncPath: settings.FileForSyncPath,
		dirForMediaPath: settings.DirForMediaPath,
	}

	return &newService
}

var (
	ErrIsNotMediaFile = errors.New("it is not a media file")
)

func (s *Service) Add(ctx context.Context, newKB NewKB) (*KB, error) {
	err := newKB.validate()
	if err != nil {
		return nil, NewDataError(fmt.Sprintf("the given values are not valid: %s", err))
	}

	kb := newKB.toKB()

	_, err = s.storage.Create(ctx, kb)
	if err != nil && errors.As(err, &ClientError{}) {
		return nil, NewDataError(fmt.Sprintf("unable to add kb due to given data: %s", err))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to add kb: %w", err)
	}

	return &kb, nil
}

func (s *Service) Update(ctx context.Context, kb KB) error {
	err := kb.validate()
	if err != nil {
		return NewDataError(fmt.Sprintf("the given values are not valid: %s", err))
	}

	err = s.storage.Update(ctx, &kb)
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
		kb := newKB.toKB()
		_, err := s.storage.Create(ctx, kb)
		if err != nil {
			result.FailedKeys[newKB.Key] = err.Error()
			continue
		}

		result.NewIDs[newKB.Key] = kb.ID
	}

	return &result, nil
}

func (s *Service) GetAllKBs(ctx context.Context, filter KBQueryFilter) (*GetAllResult, error) {
	err := filter.valid()
	if err != nil {
		return nil, fmt.Errorf("unable to get all kbs, invalid filter values: %w", err)
	}

	result, err := s.storage.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to export kbs: %w", err)
	}

	return result, nil
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

func (s *Service) SaveMedia(ctx context.Context, newKB NewKB) error {
	isNotMediaFile, err := isNotMediaFile(newKB.Value)
	if err != nil {
		return fmt.Errorf("unable to save media: %w", err)
	}

	if isNotMediaFile {
		return ErrIsNotMediaFile
	}

	if !isWebURL(newKB.Value) {
		return nil
	}

	// check if media folder exists in kbkitt config directory
	mediaFolderExist, err := filesystems.FolderExist(s.dirForMediaPath)
	if err != nil {
		return fmt.Errorf("unable to save media: %w", err)
	}

	if !mediaFolderExist {
		err = filesystems.MakeFolder(s.dirForMediaPath)
		if err != nil {
			return fmt.Errorf("unable to save media: %w", err)
		}
	}

	content, err := webs.GetWebMediaFile(newKB.Value)
	if err != nil {
		return fmt.Errorf("unable to save media: %w", err)
	}

	mediaType := newKB.MediaType
	fmt.Println("content-type", http.DetectContentType(content))

	mediaFileName := filepath.Join(s.dirForMediaPath, newKB.Key+"."+mediaType)

	err = filesystems.SaveFile(mediaFileName, content)
	if err != nil {
		return fmt.Errorf("unable to save media: %w", err)
	}

	return nil
}

func isNotMediaFile(urlpath string) (bool, error) {
	if isWebURL(urlpath) {
		return false, nil
	}

	info, err := filesystems.CheckFile(urlpath)
	if err != nil {
		return true, fmt.Errorf("unable to verify if the media exists: %w", err)
	}

	if info == nil { // not exist
		return true, nil
	}

	if info.IsDir { // it is a dir so it is not media file
		return true, nil
	}

	return false, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*KB, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("the given id is not valid, because it is empty")
	}

	kb, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get kb: %w", err)
	}

	return kb, nil
}

func (s *Service) GetByKey(ctx context.Context, key string) (*KB, error) {
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("the given key is not valid, because it is empty")
	}

	kb, err := s.storage.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get kb: %w", err)
	}

	return kb, nil
}

func (s *Service) Search(ctx context.Context, filter KBQueryFilter) (*SearchResult, error) {
	if filter.nothingToLookFor() {
		return nil, nil
	}

	kbs, err := s.storage.Search(ctx, filter)
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
