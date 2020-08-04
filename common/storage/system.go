package storage

import "encoding/binary"

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
