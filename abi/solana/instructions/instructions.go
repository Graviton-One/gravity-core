package instructions

import (
	"encoding/hex"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

func UpdateConsulsInstruction(fromAccount, programData, targetProgramID, multisigId common.PublicKey, signers []common.PublicKey, Bft uint8, Round uint64, Consuls []common.PublicKey) types.Instruction {
	consuls := []byte{}
	for i := 0; i < int(Bft); i++ {
		// acc := types.NewAccount()
		// Consuls[i][1:]
		consuls = append(consuls, Consuls[i][:]...)
	}
	data, err := common.SerializeData(struct {
		Instruction uint8
		Bft         uint8
		Round       uint64
		Consuls     []byte
	}{
		Instruction: 1,
		Bft:         Bft,
		Round:       Round,
		Consuls:     consuls,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("--------- RAW INSTRUCTION DATA -----------")
	fmt.Printf("%s\n", hex.EncodeToString(data))
	fmt.Println("------- END RAW INSTRUCTION DATA ---------")

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: true},
		{PubKey: programData, IsSigner: false, IsWritable: true},
		{PubKey: multisigId, IsSigner: false, IsWritable: true},
	}
	for _, s := range signers {
		accounts = append(accounts, types.AccountMeta{PubKey: s, IsSigner: true, IsWritable: false})
	}
	return types.Instruction{
		Accounts:  accounts,
		ProgramID: targetProgramID,
		Data:      data,
	}
}

func NebulaUpdateOraclesInstruction(fromAccount, targetProgramID, nebulaDataAccount, multisigAccount common.PublicKey, signers []common.PublicKey, Round uint64, Oracles []common.PublicKey, Bft uint8) types.Instruction {
	/*
			UpdateOracles {
		        new_oracles: Vec<Pubkey>,
		        new_round: PulseID,
		    }
	*/
	oracles := []byte{}
	for i := 0; i < len(Oracles); i++ {
		oracles = append(oracles, Oracles[i][:]...)
	}
	data, err := common.SerializeData(struct {
		Instruction uint8
		Bft         uint8
		NewOracles  []byte
		Round       uint64
	}{
		Instruction: 1,
		Bft:         Bft,
		NewOracles:  oracles,
		Round:       Round,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("--------- RAW INSTRUCTION DATA -----------")
	fmt.Printf("%s\n", hex.EncodeToString(data))
	fmt.Println("------- END RAW INSTRUCTION DATA ---------")

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: true},
		{PubKey: nebulaDataAccount, IsSigner: false, IsWritable: true},
		{PubKey: multisigAccount, IsSigner: false, IsWritable: true},
	}
	for _, s := range signers {
		accounts = append(accounts, types.AccountMeta{PubKey: s, IsSigner: true, IsWritable: false})
	}
	return types.Instruction{
		Accounts:  accounts,
		ProgramID: targetProgramID,
		Data:      data,
	}
}

func NebulaAddPulseInstruction(fromAccount, targetProgramID, nebulaId, multisigId, nebulaState common.PublicKey, signers []common.PublicKey, PulseId uint64, hash []byte) types.Instruction {
	/*
			SendHashValue {
		        data_hash: Vec<u8>,
		    }
	*/
	newHash := [64]byte{}
	copy(newHash[:], hash)

	data, err := common.SerializeData(struct {
		Instruction uint8
		//PulseID     uint64
		Hash []byte
	}{
		Instruction: 2,
		//PulseID:     PulseId,
		Hash: newHash[:],
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("--------- RAW INSTRUCTION DATA -----------")
	fmt.Printf("%s\n", hex.EncodeToString(data))
	fmt.Println("------- END RAW INSTRUCTION DATA ---------")

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: true},
		{PubKey: nebulaState, IsSigner: false, IsWritable: true},
		{PubKey: multisigId, IsSigner: false, IsWritable: true},
	}
	for _, s := range signers {
		accounts = append(accounts, types.AccountMeta{PubKey: s, IsSigner: true, IsWritable: false})
	}
	return types.Instruction{
		Accounts:  accounts,
		ProgramID: targetProgramID,
		Data:      data,
	}
}

func NebulaSendValueToSubsInstruction(fromAccount,
	targetProgramID, nebulaId, nebulaState, nebulaMultisig common.PublicKey,
	ibportProgramAccount, ibportDataAccount, tokenProgramAddress, recipient, ibPortPDA common.PublicKey,
	DataType uint8, value []byte, PulseId uint64, SubscriptionID [16]byte) types.Instruction {
	/*
			SendValueToSubs {
		        data_type: DataType,
		        pulse_id: PulseID,
		        subscription_id: SubscriptionID,
				value: Vec<u8>
		    },
	*/
	data, err := common.SerializeData(struct {
		Instruction    uint8
		Value          []byte
		DataType       uint8
		PulseID        uint64
		SubscriptionID [16]byte
	}{
		Instruction:    3,
		Value:          value,
		DataType:       DataType,
		PulseID:        PulseId,
		SubscriptionID: SubscriptionID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("--------- RAW INSTRUCTION DATA -----------")
	fmt.Printf("%s\n", hex.EncodeToString(data))
	fmt.Println("------- END RAW INSTRUCTION DATA ---------")

	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: true},
		{PubKey: nebulaState, IsSigner: false, IsWritable: true},
		{PubKey: nebulaMultisig, IsSigner: false, IsWritable: true},

		{PubKey: common.TokenProgramID, IsWritable: false, IsSigner: false},
		{PubKey: ibportProgramAccount, IsWritable: false, IsSigner: false},
		{PubKey: ibportDataAccount, IsWritable: true, IsSigner: false},
		{PubKey: tokenProgramAddress, IsWritable: true, IsSigner: false},
		{PubKey: recipient, IsWritable: true, IsSigner: false},
		{PubKey: ibPortPDA, IsWritable: false, IsSigner: false},
	}
	return types.Instruction{
		Accounts:  accounts,
		ProgramID: targetProgramID,
		Data:      data,
	}
}
