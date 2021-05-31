package instructions

import (
	"bytes"
	"errors"
	"math"

	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

// `NebulaSerializer` is a partial implementation of the `Serializer` interface.
// It is used as an embedded struct by the Bincode and BCS serializers.
type NebulaSerializer struct {
	Buffer               bytes.Buffer
	containerDepthBudget uint64
}

func NewNebulaSerializer(max_container_depth uint64) *NebulaSerializer {
	s := new(NebulaSerializer)
	s.containerDepthBudget = max_container_depth
	return s
}

func (d *NebulaSerializer) IncreaseContainerDepth() error {
	if d.containerDepthBudget == 0 {
		return errors.New("exceeded maximum container depth")
	}
	d.containerDepthBudget -= 1
	return nil
}

func (d *NebulaSerializer) DecreaseContainerDepth() {
	d.containerDepthBudget += 1
}

// `serializeLen` to be provided by the extending struct.
// func (s *NebulaSerializer) SerializeBytes(value []byte, serializeLen func(uint64) error) error {
// 	serializeLen(uint64(len(value)))
// 	s.Buffer.Write(value)
// 	return nil
// }
// `serializeLen` to be provided by the extending struct.
func (s *NebulaSerializer) SerializeBytes(value []byte) error {
	//serializeLen(uint64(len(value)))
	s.Buffer.Write(value)
	return nil
}

// `serializeLen` to be provided by the extending struct.
func (s *NebulaSerializer) SerializeStr(value string) error {
	return s.SerializeBytes([]byte(value))
}

func (s *NebulaSerializer) SerializeBool(value bool) error {
	if value {
		return s.Buffer.WriteByte(1)
	}
	return s.Buffer.WriteByte(0)
}

func (s *NebulaSerializer) SerializeUnit(value struct{}) error {
	return nil
}

// SerializeChar is unimplemented.
func (s *NebulaSerializer) SerializeChar(value rune) error {
	return errors.New("unimplemented")
}

func (s *NebulaSerializer) SerializeU8(value uint8) error {
	s.Buffer.WriteByte(byte(value))
	return nil
}

func (s *NebulaSerializer) SerializeU16(value uint16) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	return nil
}

func (s *NebulaSerializer) SerializeU32(value uint32) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	s.Buffer.WriteByte(byte(value >> 16))
	s.Buffer.WriteByte(byte(value >> 24))
	return nil
}

func (s *NebulaSerializer) SerializeU64(value uint64) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	s.Buffer.WriteByte(byte(value >> 16))
	s.Buffer.WriteByte(byte(value >> 24))
	s.Buffer.WriteByte(byte(value >> 32))
	s.Buffer.WriteByte(byte(value >> 40))
	s.Buffer.WriteByte(byte(value >> 48))
	s.Buffer.WriteByte(byte(value >> 56))
	return nil
}

func (s *NebulaSerializer) SerializeU128(value serde.Uint128) error {
	s.SerializeU64(value.Low)
	s.SerializeU64(value.High)
	return nil
}

func (s *NebulaSerializer) SerializeI8(value int8) error {
	s.SerializeU8(uint8(value))
	return nil
}

func (s *NebulaSerializer) SerializeI16(value int16) error {
	s.SerializeU16(uint16(value))
	return nil
}

func (s *NebulaSerializer) SerializeI32(value int32) error {
	s.SerializeU32(uint32(value))
	return nil
}

func (s *NebulaSerializer) SerializeI64(value int64) error {
	s.SerializeU64(uint64(value))
	return nil
}

func (s *NebulaSerializer) SerializeI128(value serde.Int128) error {
	s.SerializeU64(value.Low)
	s.SerializeI64(value.High)
	return nil
}

func (s *NebulaSerializer) SerializeF32(value float32) error {
	n := math.Float32bits(value)
	s.Buffer.WriteByte(byte(n))
	s.Buffer.WriteByte(byte(n >> 8))
	s.Buffer.WriteByte(byte(n >> 16))
	s.Buffer.WriteByte(byte(n >> 24))
	return nil
}

func (s *NebulaSerializer) SerializeF64(value float64) error {
	n := math.Float64bits(value)
	s.Buffer.WriteByte(byte(n))
	s.Buffer.WriteByte(byte(n >> 8))
	s.Buffer.WriteByte(byte(n >> 16))
	s.Buffer.WriteByte(byte(n >> 24))
	s.Buffer.WriteByte(byte(n >> 32))
	s.Buffer.WriteByte(byte(n >> 40))
	s.Buffer.WriteByte(byte(n >> 48))
	s.Buffer.WriteByte(byte(n >> 56))
	return nil
}

func (s *NebulaSerializer) SerializeOptionTag(value bool) error {
	return s.SerializeBool(value)
}

func (s *NebulaSerializer) GetBufferOffset() uint64 {
	return uint64(s.Buffer.Len())
}

func (s *NebulaSerializer) GetBytes() []byte {
	return s.Buffer.Bytes()
}

func (s *NebulaSerializer) SerializeLen(value uint64) error {
	return s.SerializeU64(value)
}
func (s *NebulaSerializer) SerializeVariantIndex(value uint32) error {
	return s.SerializeU32(value)
}

func (s *NebulaSerializer) SortMapEntries(offsets []uint64) {

}
