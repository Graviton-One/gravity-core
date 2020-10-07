package storage

import (
	"encoding/binary"
	"fmt"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/config"
)

func (storage *Storage) LastHeight() (uint64, error) {
	b, err := storage.getValue([]byte(LastHeightKey))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}

func (storage *Storage) SetLastHeight(height uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], height)
	err := storage.setValue([]byte(LastHeightKey), b[:])
	if err != nil {
		return err
	}

	return err
}

func (storage *Storage) ConsulsCount() (int, error) {
	b, err := storage.getValue([]byte(ConsulsCountKey))
	if err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint64(b)), nil
}

func (storage *Storage) SetConsulsCount(consulsCount int) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(consulsCount))
	err := storage.setValue([]byte(ConsulsCountKey), b[:])
	if err != nil {
		return err
	}

	return err
}

func (storage *Storage) SetLastRoundApproved(roundId uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], roundId)
	err := storage.setValue([]byte(LastRoundApproved), b[:])
	if err != nil {
		return err
	}

	return err
}
func (storage *Storage) LastRoundApproved() (uint64, error) {
	b, err := storage.getValue([]byte(LastRoundApproved))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}


func (storage *Storage) Validators() (*[]account.ConsulPubKey, error) {
	consulScores, err := storage.Scores()

	if err != nil {
		return nil, err
	}

	keys := make([]account.ConsulPubKey, len(consulScores))

	for consulPubKey, _ := range consulScores {
		keys = append(keys, consulPubKey)
	}

	return &keys, nil
}

func (storage *Storage) ValidatorDetails() (*config.ValidatorDetails, error) {
	// Read current config (on every call)
	if storage.AppDetailsDelegate == nil {
		return nil, fmt.Errorf("no details provided")
	}

	return storage.AppDetailsDelegate, nil
}