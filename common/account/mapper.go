package account

import (
	"fmt"
)

type Mapper struct {
	strToByte map[string]byte
	byteToStr map[byte]string
	types     map[byte]string
}

var ChainMapper Mapper

func (m *Mapper) Assign(arg map[string]ChainType) {
	m.strToByte = make(map[string]byte)
	m.byteToStr = make(map[byte]string)
	m.types = make(map[byte]string)
	for k, v := range arg {
		m.byteToStr[byte(v)] = k
		m.strToByte[k] = byte(v)
	}
}
func (m *Mapper) ApendAdaptor(id byte, chaintype string) {
	m.types[id] = chaintype
}
func (m *Mapper) ToType(id byte) (byte, error) {
	v, ok := m.types[id]
	if ok {
		v2, ok2 := m.strToByte[v]
		if ok2 {
			return v2, nil
		}
		return 0, fmt.Errorf("Chain not found")
	}
	return 0, fmt.Errorf("Chain not found")
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
