package instructions

import (
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

func UpdateConsulsInstruction(fromAccount, programData, targetProgramID common.PublicKey, signers []common.PublicKey, Bft uint8, Round uint64, Consuls []common.PublicKey) types.Instruction {
	consuls := []byte{}
	for i := 0; i < 3; i++ {
		acc := types.NewAccount()
		consuls = append(consuls, acc.PublicKey.Bytes()...)
	}
	data, err := common.SerializeData(struct {
		Instruction uint8
		Bft         uint8
		Consuls     []byte
		Round       uint64
	}{
		Instruction: 0,
		Bft:         3,
		Round:       Round,
		Consuls:     consuls,
	})
	if err != nil {
		panic(err)
	}
	accounts := []types.AccountMeta{
		{PubKey: fromAccount, IsSigner: true, IsWritable: true},
		{PubKey: programData, IsSigner: false, IsWritable: true},
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
