package adaptors

import (
	"context"
	"crypto/ed25519"
	"sort"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"go.uber.org/zap"

	solana "github.com/portto/solana-go-sdk/client"
	solana_common "github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/types"
)

type SortablePubkey []solana_common.PublicKey

func (spk SortablePubkey) Len() int { return len(spk) }
func (spk SortablePubkey) Less(i, j int) bool {
	for x := range spk[i] {
		if spk[i][x] == spk[j][x] {
			continue
		}
		return spk[i][x] < spk[j][x]
	}
	return false
}

func (spk SortablePubkey) Swap(i, j int) { spk[i], spk[j] = spk[j], spk[i] }
func (spk SortablePubkey) ToPubKeys() []solana_common.PublicKey {
	res := []solana_common.PublicKey{}
	for _, v := range spk {
		res = append(res, v)
	}
	return res
}

type SolanaAdapter struct {
	account types.Account
	client  *solana.Client
}

func NewSolanaAdaptor(privKey []byte, nodeUrl string, ctx context.Context) (*SolanaAdapter, error) {

	account := types.AccountFromPrivateKeyBytes(privKey)
	solClient := solana.NewClient(nodeUrl)
	adapter := SolanaAdapter{
		client:  solClient,
		account: account,
	}

	return &adapter, nil
}

func (s *SolanaAdapter) GetHeight(ctx context.Context) (uint64, error) {
	info, err := s.client.GetEpochInfo(solana.CommitmentFinalized)
	if err != nil {
		return 0, err
	}
	return uint64(info.BlockHeight), nil
}

func (s *SolanaAdapter) WaitTx(id string, ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.account.PrivateKey, msg), nil
}

func (s *SolanaAdapter) PubKey() account.OraclesPubKey {
	var pubKey account.OraclesPubKey
	copy(pubKey[:], append([]byte{0}, s.account.PublicKey.Bytes()[0:32]...))
	return pubKey
}

func (s *SolanaAdapter) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64) ([]byte, error) {

	consuls := SortablePubkey{}
	for _, v := range consulsAddresses {
		if v == nil {
			consuls = append(consuls, types.Account{}.PublicKey)
			continue
		}
		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		consuls = append(consuls, pubKey)
	}
	sort.Sort(&SortablePubkey{})
	solanaConsuls := consuls.ToPubKeys()

	res, err := s.client.GetRecentBlockhash()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return nil, err
	}
	message := types.NewMessage(
		account.PublicKey,
		[]types.Instruction{
			NewInitGravityContractInstruction(
				account.PublicKey, dataAcc, multisigAcc, program, 3, 1, [5][32]byte{},
			),
		},
		res.Blockhash,
	)
	// hash, err := adaptor.gravityContract.HashNewConsuls(nil, oraclesAddresses, big.NewInt(roundId))
	// if err != nil {
	// 	return nil, err
	// }

	// sign, err := adaptor.Sign(hash[:])
	// if err != nil {
	// 	return nil, err
	// }

	// return sign, nil

	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastRound(ctx context.Context) (uint64, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	panic("not implemented") // TODO: Implement
}
