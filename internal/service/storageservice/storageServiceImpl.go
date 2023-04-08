package storageservice

import (
	"context"

	"github.com/rs/zerolog/log"

	"protected-storage-server/internal/repositories/rawdatarepository"
)

var _ StorageService = &storageServiceImpl{}

type storageServiceImpl struct {
	rawDataRepository rawdatarepository.RawDataRepository
}

// SaveRawData метод для сохранения произвольных текстовых данных
func (s storageServiceImpl) SaveRawData(ctx context.Context, name, data, userID string) error {
	log.Info().Msgf("storageservice: save raw data for user with ID %s", userID)
	return s.rawDataRepository.Save(ctx, userID, name, data)
}

// GetRawData метод для получения произвольных текстовых данных
func (s storageServiceImpl) GetRawData(ctx context.Context, name, userID string) (string, error) {
	log.Info().Msgf("storageservice: get raw data with name %s for user with ID %s", name, userID)
	return s.rawDataRepository.GetByName(ctx, userID, name)
}

// New конструктор UserService
func New(rawDataRepository rawdatarepository.RawDataRepository) StorageService {
	return &storageServiceImpl{
		rawDataRepository,
	}
}
