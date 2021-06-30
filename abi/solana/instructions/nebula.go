package instructions

import (
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

// type Pulse struct {
// 	DataHash []uint8
// 	Height   uint64
// }
// type SubscriptionID [16]uint8

// // struct Pulse {
// //     pub data_hash: Vec<u8>,
// //     pub height: u128,
// // }
// /*
// pub struct NebulaContract {
//     pub rounds_dict: HashMap<PulseID, bool>,
//     subscriptions_queue: NebulaQueue<SubscriptionID>,
//     pub oracles: Vec<Pubkey>,

//     pub bft: u8,
//     pub multisig_account: Pubkey,
//     pub gravity_contract: Pubkey,
//     pub data_type: DataType,
//     pub last_round: PulseID,

//     // subscription_ids: Vec<SubscriptionID>,
//     pub last_pulse_id: PulseID,

//     subscriptions_map: HashMap<SubscriptionID, Subscription>,
//     pulses_map: HashMap<PulseID, Pulse>,
//     is_pulse_sent: HashMap<PulseID, HashMap<SubscriptionID, bool>>,

//     pub is_initialized: bool,
//     pub initializer_pubkey: Pubkey,
// }
// */

// // serde: deserialize,serialize
// type NebulaContract struct {
// 	RoundsDict         map[uint64]bool
// 	SubscriptionsQueue []SubscriptionID
// 	Oracles            [][32]byte
// 	Bft                uint8
// 	MultisigAccount    [32]byte
// 	GravityContract    [32]byte
// 	DataType           uint8
// 	LastRound          uint64
// 	LastPulseID        uint64
// 	SubscriptionsMap   map[SubscriptionID]Subscription
// 	PulsesMap          map[uint64]Pulse
// 	IsPulseSent        map[uint64]map[SubscriptionID]bool
// 	IsInitialized      bool
// 	InitializerPubKey  [32]byte
// }

// type Subscription struct {
// 	Sender           [32]byte
// 	ContractAddress  [32]byte
// 	MinConfirmations uint8
// 	Reward           uint64
// }

func InitNebulaInstruction(fromAccount, gravityProgramData, programID, nebulaAccount common.PublicKey, Bft uint8, DataType uint8, Oracles []common.PublicKey) (*types.Instruction, error) {
	/*
	   InitContract {
	       nebula_data_type: DataType,
	       gravity_contract_program_id: Pubkey,
	       initial_oracles: Vec<Pubkey>,
	       oracles_bft: u8,
	   },
	*/
	data, err := common.SerializeData(struct {
		Instruction        uint8
		DataType           uint8
		GravityProgramData common.PublicKey
		InitialOracles     []common.PublicKey
		Bft                uint8
	}{
		Instruction:        0,
		DataType:           DataType,
		GravityProgramData: gravityProgramData,
		InitialOracles:     Oracles,
		Bft:                Bft,
	})
	if err != nil {
		return nil, err
	}

	return &types.Instruction{
		Accounts: []types.AccountMeta{
			{PubKey: fromAccount, IsSigner: true, IsWritable: false},
			{PubKey: nebulaAccount, IsSigner: false, IsWritable: true},
			{PubKey: gravityProgramData, IsSigner: false, IsWritable: true},
		},
		ProgramID: programID,
		Data:      data,
	}, nil
}

func UpdateOraclesInstruction(fromAccount, programID, nebulaAccount common.PublicKey, CurrentOracles, NewOracles []common.PublicKey, pulseID uint64) (*types.Instruction, error) {
	/*UpdateOracles {
	  	new_oracles: Vec<Pubkey>,
	  	new_round: PulseID,
	  },
	*/
	data, err := common.SerializeData(struct {
		Instruction uint8
		Bft         uint8
		NewOracles  []common.PublicKey
		PulseID     uint64
	}{
		Instruction: 1,
		Bft:         uint8(len(CurrentOracles)),
		NewOracles:  NewOracles,
		PulseID:     uint64(pulseID),
	})
	if err != nil {
		return nil, err
	}

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: false},
		{PubKey: nebulaAccount, IsSigner: false, IsWritable: true},
	}
	for _, oracle := range CurrentOracles {
		accounts = append(accounts, types.AccountMeta{
			PubKey:     oracle,
			IsSigner:   true,
			IsWritable: false,
		})
	}
	return &types.Instruction{
		Accounts:  accounts,
		ProgramID: programID,
		Data:      data,
	}, nil

}

func SendHashValueInstruction(fromAccount, programID, nebulaAccount common.PublicKey, currentOracles []common.PublicKey, dataHash [16]byte) (*types.Instruction, error) {
	/*
	   SendHashValue {
	   	data_hash: UUID,
	   },
	*/
	data, err := common.SerializeData(struct {
		Instruction uint8
		DataHash    [16]byte
	}{
		Instruction: 2,
		DataHash:    dataHash,
	})
	if err != nil {
		return nil, err
	}

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: false},
		{PubKey: nebulaAccount, IsSigner: false, IsWritable: true},
	}
	for _, oracle := range currentOracles {
		accounts = append(accounts, types.AccountMeta{
			PubKey:     oracle,
			IsSigner:   true,
			IsWritable: false,
		})
	}
	return &types.Instruction{
		Accounts:  accounts,
		ProgramID: programID,
		Data:      data,
	}, nil
}

func SendValueToSubsInstruction(fromAccount, programID, nebulaAccount common.PublicKey, currentOracles []common.PublicKey, pulseID uint64, subscriptionID [16]byte, dType DataType, Hash []byte) (*types.Instruction, error) {

	/*
	   SendValueToSubs {
	   	data_type: DataType,
	   	pulse_id: PulseID,
	   	subscription_id: UUID,
	   },
	*/
	data, err := common.SerializeData(struct {
		Instruction   uint8
		DataHash      []byte
		DataType      DataType
		PulseID       uint64
		SubsriptionID [16]byte
	}{
		Instruction:   3,
		DataHash:      Hash,
		DataType:      dType,
		PulseID:       pulseID,
		SubsriptionID: subscriptionID,
	})
	if err != nil {
		return nil, err
	}

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: false},
		{PubKey: nebulaAccount, IsSigner: false, IsWritable: true},
	}
	for _, oracle := range currentOracles {
		accounts = append(accounts, types.AccountMeta{
			PubKey:     oracle,
			IsSigner:   false,
			IsWritable: false,
		})
	}
	return &types.Instruction{
		Accounts:  accounts,
		ProgramID: programID,
		Data:      data,
	}, nil

}
