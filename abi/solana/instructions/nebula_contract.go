package instructions

import (
	"fmt"

	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

type DataType interface {
	isDataType()
	Serialize(serializer serde.Serializer) error
	Value() uint8
}

func DeserializeDataType(deserializer serde.Deserializer) (DataType, error) {
	index, err := deserializer.DeserializeVariantIndex()
	if err != nil {
		return nil, err
	}

	switch index {
	case 0:
		if val, err := load_DataType__Int64(deserializer); err == nil {
			return &val, nil
		} else {
			return nil, err
		}

	case 1:
		if val, err := load_DataType__String(deserializer); err == nil {
			return &val, nil
		} else {
			return nil, err
		}

	case 2:
		if val, err := load_DataType__Bytes(deserializer); err == nil {
			return &val, nil
		} else {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Unknown variant index for DataType: %d", index)
	}
}

type DataType__Int64 struct {
}

func (*DataType__Int64) isDataType() {}
func (*DataType__Int64) Value() uint8 {
	return 0
}
func (obj *DataType__Int64) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(0)
	serializer.DecreaseContainerDepth()
	return nil
}

func load_DataType__Int64(deserializer serde.Deserializer) (DataType__Int64, error) {
	var obj DataType__Int64
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}

type DataType__String struct {
}

func (*DataType__String) isDataType() {}

func (obj *DataType__String) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(1)
	serializer.DecreaseContainerDepth()
	return nil
}
func (*DataType__String) Value() uint8 {
	return 1
}
func load_DataType__String(deserializer serde.Deserializer) (DataType__String, error) {
	var obj DataType__String
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}

type DataType__Bytes struct {
}

func (*DataType__Bytes) Value() uint8 {
	return 2
}
func (*DataType__Bytes) isDataType() {}

func (obj *DataType__Bytes) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	serializer.SerializeVariantIndex(2)
	serializer.DecreaseContainerDepth()
	return nil
}

func load_DataType__Bytes(deserializer serde.Deserializer) (DataType__Bytes, error) {
	var obj DataType__Bytes
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}

type NebulaContract struct {
	RoundsDict         map[uint64]bool
	SubscriptionsQueue [][16]uint8
	Oracles            []Pubkey
	Bft                uint8
	MultisigAccount    Pubkey
	GravityContract    Pubkey
	DataType           DataType
	LastRound          uint64
	LastPulseId        uint64
	SubscriptionsMap   map[[16]uint8]Subscription
	PulsesMap          map[uint64]Pulse
	IsPulseSent        map[uint64]bool
	IsInitialized      bool
	InitializerPubkey  Pubkey
}

func NewNebulaContract() *NebulaContract {
	c := NebulaContract{}
	c.RoundsDict = make(map[uint64]bool)
	c.SubscriptionsQueue = make([][16]uint8, 0)
	c.Oracles = make([]Pubkey, 0)
	c.MultisigAccount = Pubkey{}
	c.GravityContract = Pubkey{}
	c.DataType = &DataType__Bytes{}
	c.SubscriptionsMap = make(map[[16]uint8]Subscription)
	c.PulsesMap = make(map[uint64]Pulse)
	c.IsPulseSent = make(map[uint64]bool)
	c.InitializerPubkey = Pubkey{}
	return &c
}

func (obj *NebulaContract) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serialize_map_u64_to_bool(obj.RoundsDict, serializer); err != nil {
		return err
	}
	if err := serialize_vector_array16_u8_array(obj.SubscriptionsQueue, serializer); err != nil {
		return err
	}
	if err := serialize_vector_Pubkey(obj.Oracles, serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU8(obj.Bft); err != nil {
		return err
	}
	if err := obj.MultisigAccount.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.GravityContract.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.DataType.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.LastRound); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.LastPulseId); err != nil {
		return err
	}
	if err := serialize_map_array16_u8_array_to_Subscription(obj.SubscriptionsMap, serializer); err != nil {
		return err
	}
	if err := serialize_map_u64_to_Pulse(obj.PulsesMap, serializer); err != nil {
		return err
	}
	if err := serialize_map_u64_to_bool(obj.IsPulseSent, serializer); err != nil {
		return err
	}
	if err := serializer.SerializeBool(obj.IsInitialized); err != nil {
		return err
	}
	if err := obj.InitializerPubkey.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func DeserializeNebulaContract(deserializer serde.Deserializer) (NebulaContract, error) {
	var obj NebulaContract
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	if val, err := deserialize_map_u64_to_bool(deserializer); err == nil {
		obj.RoundsDict = val
	} else {
		return obj, err
	}
	if val, err := deserialize_vector_array16_u8_array(deserializer); err == nil {
		obj.SubscriptionsQueue = val
	} else {
		return obj, err
	}
	if val, err := deserialize_vector_Pubkey(deserializer); err == nil {
		obj.Oracles = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU8(); err == nil {
		obj.Bft = val
	} else {
		return obj, err
	}
	if val, err := DeserializePubkey(deserializer); err == nil {
		obj.MultisigAccount = val
	} else {
		return obj, err
	}
	if val, err := DeserializePubkey(deserializer); err == nil {
		obj.GravityContract = val
	} else {
		return obj, err
	}
	if val, err := DeserializeDataType(deserializer); err == nil {
		obj.DataType = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU64(); err == nil {
		obj.LastRound = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU64(); err == nil {
		obj.LastPulseId = val
	} else {
		return obj, err
	}
	if val, err := deserialize_map_array16_u8_array_to_Subscription(deserializer); err == nil {
		obj.SubscriptionsMap = val
	} else {
		return obj, err
	}
	if val, err := deserialize_map_u64_to_Pulse(deserializer); err == nil {
		obj.PulsesMap = val
	} else {
		return obj, err
	}
	if val, err := deserialize_map_u64_to_bool(deserializer); err == nil {
		obj.IsPulseSent = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeBool(); err == nil {
		obj.IsInitialized = val
	} else {
		return obj, err
	}
	if val, err := DeserializePubkey(deserializer); err == nil {
		obj.InitializerPubkey = val
	} else {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}

type Pubkey [32]uint8

func (obj *Pubkey) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serialize_array32_u8_array((([32]uint8)(*obj)), serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func DeserializePubkey(deserializer serde.Deserializer) (Pubkey, error) {
	var obj [32]uint8
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return (Pubkey)(obj), err
	}
	if val, err := deserialize_array32_u8_array(deserializer); err == nil {
		obj = val
	} else {
		return ((Pubkey)(obj)), err
	}
	deserializer.DecreaseContainerDepth()
	return (Pubkey)(obj), nil
}

type Pulse struct {
	DataHash []uint8
	Height   uint64
}

func (obj *Pulse) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serialize_vector_u8(obj.DataHash, serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.Height); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func DeserializePulse(deserializer serde.Deserializer) (Pulse, error) {
	var obj Pulse
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	if val, err := deserialize_vector_u8(deserializer); err == nil {
		obj.DataHash = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU64(); err == nil {
		obj.Height = val
	} else {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}

type Subscription struct {
	Sender           Pubkey
	ContractAddress  Pubkey
	MinConfirmations uint8
	Reward           uint64
}

func (obj *Subscription) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := obj.Sender.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.ContractAddress.Serialize(serializer); err != nil {
		return err
	}
	if err := serializer.SerializeU8(obj.MinConfirmations); err != nil {
		return err
	}
	if err := serializer.SerializeU64(obj.Reward); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func DeserializeSubscription(deserializer serde.Deserializer) (Subscription, error) {
	var obj Subscription
	if err := deserializer.IncreaseContainerDepth(); err != nil {
		return obj, err
	}
	if val, err := DeserializePubkey(deserializer); err == nil {
		obj.Sender = val
	} else {
		return obj, err
	}
	if val, err := DeserializePubkey(deserializer); err == nil {
		obj.ContractAddress = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU8(); err == nil {
		obj.MinConfirmations = val
	} else {
		return obj, err
	}
	if val, err := deserializer.DeserializeU64(); err == nil {
		obj.Reward = val
	} else {
		return obj, err
	}
	deserializer.DecreaseContainerDepth()
	return obj, nil
}
func serialize_array16_u8_array(value [16]uint8, serializer serde.Serializer) error {
	for _, item := range value {
		if err := serializer.SerializeU8(item); err != nil {
			return err
		}
	}
	return nil
}

func deserialize_array16_u8_array(deserializer serde.Deserializer) ([16]uint8, error) {
	var obj [16]uint8
	for i := range obj {
		if val, err := deserializer.DeserializeU8(); err == nil {
			obj[i] = val
		} else {
			return obj, err
		}
	}
	return obj, nil
}

func serialize_array32_u8_array(value [32]uint8, serializer serde.Serializer) error {
	for _, item := range value {
		if err := serializer.SerializeU8(item); err != nil {
			return err
		}
	}
	return nil
}

func deserialize_array32_u8_array(deserializer serde.Deserializer) ([32]uint8, error) {
	var obj [32]uint8
	for i := range obj {
		if val, err := deserializer.DeserializeU8(); err == nil {
			obj[i] = val
		} else {
			return obj, err
		}
	}
	return obj, nil
}

func serialize_map_array16_u8_array_to_Subscription(value map[[16]uint8]Subscription, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	offsets := make([]uint64, len(value))
	count := 0
	for k, v := range value {
		offsets[count] = serializer.GetBufferOffset()
		count += 1
		if err := serialize_array16_u8_array(k, serializer); err != nil {
			return err
		}
		if err := v.Serialize(serializer); err != nil {
			return err
		}
	}
	serializer.SortMapEntries(offsets)
	return nil
}

func deserialize_map_array16_u8_array_to_Subscription(deserializer serde.Deserializer) (map[[16]uint8]Subscription, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make(map[[16]uint8]Subscription)
	previous_slice := serde.Slice{0, 0}
	for i := 0; i < int(length); i++ {
		var slice serde.Slice
		slice.Start = deserializer.GetBufferOffset()
		var key [16]uint8
		if val, err := deserialize_array16_u8_array(deserializer); err == nil {
			key = val
		} else {
			return nil, err
		}
		slice.End = deserializer.GetBufferOffset()
		if i > 0 {
			err := deserializer.CheckThatKeySlicesAreIncreasing(previous_slice, slice)
			if err != nil {
				return nil, err
			}
		}
		previous_slice = slice
		if val, err := DeserializeSubscription(deserializer); err == nil {
			obj[key] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}

func serialize_map_u64_to_Pulse(value map[uint64]Pulse, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	offsets := make([]uint64, len(value))
	count := 0
	for k, v := range value {
		offsets[count] = serializer.GetBufferOffset()
		count += 1
		if err := serializer.SerializeU64(k); err != nil {
			return err
		}
		if err := v.Serialize(serializer); err != nil {
			return err
		}
	}
	serializer.SortMapEntries(offsets)
	return nil
}

func deserialize_map_u64_to_Pulse(deserializer serde.Deserializer) (map[uint64]Pulse, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make(map[uint64]Pulse)
	previous_slice := serde.Slice{0, 0}
	for i := 0; i < int(length); i++ {
		var slice serde.Slice
		slice.Start = deserializer.GetBufferOffset()
		var key uint64
		if val, err := deserializer.DeserializeU64(); err == nil {
			key = val
		} else {
			return nil, err
		}
		slice.End = deserializer.GetBufferOffset()
		if i > 0 {
			err := deserializer.CheckThatKeySlicesAreIncreasing(previous_slice, slice)
			if err != nil {
				return nil, err
			}
		}
		previous_slice = slice
		if val, err := DeserializePulse(deserializer); err == nil {
			obj[key] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}

func serialize_map_u64_to_bool(value map[uint64]bool, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	offsets := make([]uint64, len(value))
	count := 0
	for k, v := range value {
		offsets[count] = serializer.GetBufferOffset()
		count += 1
		if err := serializer.SerializeU64(k); err != nil {
			return err
		}
		if err := serializer.SerializeBool(v); err != nil {
			return err
		}
	}
	serializer.SortMapEntries(offsets)
	return nil
}

func deserialize_map_u64_to_bool(deserializer serde.Deserializer) (map[uint64]bool, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make(map[uint64]bool)
	previous_slice := serde.Slice{0, 0}
	for i := 0; i < int(length); i++ {
		var slice serde.Slice
		slice.Start = deserializer.GetBufferOffset()
		var key uint64
		if val, err := deserializer.DeserializeU64(); err == nil {
			key = val
		} else {
			return nil, err
		}
		slice.End = deserializer.GetBufferOffset()
		if i > 0 {
			err := deserializer.CheckThatKeySlicesAreIncreasing(previous_slice, slice)
			if err != nil {
				return nil, err
			}
		}
		previous_slice = slice
		if val, err := deserializer.DeserializeBool(); err == nil {
			obj[key] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}

func serialize_vector_Pubkey(value []Pubkey, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := item.Serialize(serializer); err != nil {
			return err
		}
	}
	return nil
}

func deserialize_vector_Pubkey(deserializer serde.Deserializer) ([]Pubkey, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make([]Pubkey, length)
	for i := range obj {
		if val, err := DeserializePubkey(deserializer); err == nil {
			obj[i] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}

func serialize_vector_array16_u8_array(value [][16]uint8, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := serialize_array16_u8_array(item, serializer); err != nil {
			return err
		}
	}
	return nil
}

func deserialize_vector_array16_u8_array(deserializer serde.Deserializer) ([][16]uint8, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make([][16]uint8, length)
	for i := range obj {
		if val, err := deserialize_array16_u8_array(deserializer); err == nil {
			obj[i] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}

func serialize_vector_u8(value []uint8, serializer serde.Serializer) error {
	if err := serializer.SerializeLen(uint64(len(value))); err != nil {
		return err
	}
	for _, item := range value {
		if err := serializer.SerializeU8(item); err != nil {
			return err
		}
	}
	return nil
}

func deserialize_vector_u8(deserializer serde.Deserializer) ([]uint8, error) {
	length, err := deserializer.DeserializeLen()
	if err != nil {
		return nil, err
	}
	obj := make([]uint8, length)
	for i := range obj {
		if val, err := deserializer.DeserializeU8(); err == nil {
			obj[i] = val
		} else {
			return nil, err
		}
	}
	return obj, nil
}
