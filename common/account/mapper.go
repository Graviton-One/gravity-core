package account

import (
	"fmt"
)

type Mapper struct {
	strToByte map[string]byte
	byteToStr map[byte]string
}

var ChainMapper Mapper

func (m *Mapper) Assign(arg map[string]ChainType) {
	m.strToByte = make(map[string]byte)
	m.byteToStr = make(map[byte]string)

	for k, v := range arg {
		m.byteToStr[byte(v)] = k
		m.strToByte[k] = byte(v)
	}
}
func (m *Mapper) ToStr(k byte) (string, error) {
	v, ok := m.byteToStr[k]
	if ok {
		return v, nil
	}
	return "", fmt.Errorf("Chain not found")
}

func (m *Mapper) ToByte(k string) (byte, error) {
	v, ok := m.strToByte[k]
	if ok {
		return v, nil
	}
	return 0, fmt.Errorf("Chain not found")
}
