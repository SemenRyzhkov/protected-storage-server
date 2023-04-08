package storageservice

import (
	"context"

	"github.com/rs/zerolog/log"

	"protected-storage-server/internal/repositories/rawdatarepository"
	"protected-storage-server/internal/security"
)

var _ StorageService = &storageServiceImpl{}

type storageServiceImpl struct {
	rawDataRepository rawdatarepository.RawDataRepository
	cipherManager     *security.CipherManager
}

// SaveRawData метод для сохранения произвольных текстовых данных
func (s storageServiceImpl) SaveRawData(ctx context.Context, name, data, userID string) error {
	log.Info().Msgf("storageservice: save raw data for user with ID %s", userID)
	savedData, err := s.cipherManager.Encrypt([]byte(data))
	log.Printf("savedData %s", savedData)
	if err != nil {
		return err
	}
	return s.rawDataRepository.Save(ctx, userID, name, savedData)
}

// GetRawData метод для получения произвольных текстовых данных
func (s storageServiceImpl) GetRawData(ctx context.Context, name, userID string) (string, error) {
	log.Info().Msgf("storageservice: get raw data with name %s for user with ID %s", name, userID)
	cipherData, err := s.rawDataRepository.GetByName(ctx, userID, name)
	if err != nil {
		return "", err
	}

	decryptData, err := s.cipherManager.Decrypt(cipherData)
	if err != nil {
		return "", err
	}
	return string(decryptData), nil

}

// New конструктор UserService
func New(rawDataRepository rawdatarepository.RawDataRepository, cipherManager *security.CipherManager) StorageService {
	return &storageServiceImpl{
		rawDataRepository,
		cipherManager,
	}
}
