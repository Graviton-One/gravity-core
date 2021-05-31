package instructions

import "github.com/portto/solana-go-sdk/common"

/*
pub initializer_pubkey: Pubkey,

pub bft: u8,
pub consuls: Vec<Pubkey>,
pub last_round: u64,
pub multisig_account: Pubkey,
*/

type GravityContract struct {
	InitializerPubkey common.PublicKey
	Bft               uint8
	Consuls           []common.PublicKey
	LastRound         uint64
	MultisigAccount   common.PublicKey
}

func NewGravityContract() GravityContract {
	c := GravityContract{}
	c.Consuls = []common.PublicKey{}
	return c
}
