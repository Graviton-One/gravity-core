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
		Consuls     []byte
		Round       uint64
	}{
		Instruction: 1,
		Bft:         Bft,
		Consuls:     consuls,
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
		{PubKey: programData, IsSigner: false, IsWritable: true},
		{PubKey: multisigId, IsSigner: false, IsWritable: true},
	}
	for _, s := range signers {
		accounts = append(accounts, types.AccountMeta{PubKey: s, IsSigner: false, IsWritable: false})
	}
	return types.Instruction{
		Accounts:  accounts,
		ProgramID: targetProgramID,
		Data:      data,
	}
}
