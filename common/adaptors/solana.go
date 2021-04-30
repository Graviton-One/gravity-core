package adaptors

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"log"
	"sort"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/abi/solana/instructions"
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
	account         types.Account
	programID       solana_common.PublicKey
	gravityContract solana_common.PublicKey
	client          *solana.Client
}

func NewSolanaAdaptor(privKey []byte, nodeUrl string, custom map[string]interface{}) (*SolanaAdapter, error) {

	account := types.AccountFromPrivateKeyBytes(privKey)
	solClient := solana.NewClient(nodeUrl)
	gravityContract, ok := custom["gravity_contract"].(string)
	if !ok {
		zap.L().Error("Cannot parse gravity contract")
	}
	programID, ok := custom["program_id"].(string)
	if !ok {
		zap.L().Error("Cannot parse gravity contract")
	}
	adapter := SolanaAdapter{
		client:          solClient,
		account:         account,
		gravityContract: solana_common.PublicKeyFromString(gravityContract),
		programID:       solana_common.PublicKeyFromString(programID),
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
	msg, err := s.createUpdateConsulsMessage(newConsulsAddresses, round)
	if err != nil {
		return "", err
	}
	serializedMessage, err := msg.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}
	solsigs := make(map[solana_common.PublicKey]types.Signature)
	selfSig, err := s.Sign(serializedMessage)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}
	solsigs[s.account.PublicKey] = selfSig
	for key, sig := range signs {
		nkey := solana_common.PublicKey{}
		copy(nkey[:], key[1:33])
		solsigs[nkey] = sig
	}

	tx, err := types.CreateTransaction(msg, solsigs)
	if err != nil {
		return "", err
	}
	rawTx, err := tx.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	txSig, err := s.client.SendRawTransaction(rawTx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	log.Println("txHash:", txSig)
	return txSig, nil
}

func (s *SolanaAdapter) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64) ([]byte, error) {
	msg, err := s.createUpdateConsulsMessage(consulsAddresses, roundId)
	if err != nil {
		return nil, err
	}
	serializedMessage, err := msg.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return nil, err
	}
	sign, err := s.Sign(serializedMessage)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func (s *SolanaAdapter) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastRound(ctx context.Context) (uint64, error) {
	r, err := s.client.GetAccountInfo(s.gravityContract.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 8,
			Offset: 130,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}
	sval, ok := r.Data.(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return 0, err
	}
	val, err := base64.RawStdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}
	round := binary.BigEndian.Uint64(val)
	return round, nil
}

func (s *SolanaAdapter) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	return false, nil //default mock
	//panic("not implemented") // TODO: Implement
}

//Custom solana methods

func (s *SolanaAdapter) GetCurrentConsuls() ([]solana_common.PublicKey, error) {
	r, err := s.client.GetAccountInfo(s.gravityContract.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 96,
			Offset: 34,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return []solana_common.PublicKey{}, err
	}
	sval, ok := r.Data.(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return []solana_common.PublicKey{}, err
	}
	val, err := base64.RawStdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return []solana_common.PublicKey{}, err
	}

	sconsuls := SortablePubkey{}
	for i := 0; i < 3*32; i += 32 { // BFT=3
		pubk := solana_common.PublicKey{}
		copy(pubk[:], val[i:i+32])
		sconsuls = append(sconsuls, pubk)
	}
	sort.Sort(&sconsuls)
	return sconsuls, nil
}

func (s *SolanaAdapter) createUpdateConsulsMessage(consulsAddresses []*account.OraclesPubKey, roundId int64) (types.Message, error) {
	currentConsuls, err := s.GetCurrentConsuls()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return types.Message{}, err
	}

	consuls := SortablePubkey{}
	for _, v := range consulsAddresses {
		if v == nil {
			consuls = append(consuls, types.Account{}.PublicKey)
			continue
		}
		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		consuls = append(consuls, pubKey)
	}
	sort.Sort(&consuls)
	solanaConsuls := consuls.ToPubKeys()

	res, err := s.client.GetRecentBlockhash()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return types.Message{}, err
	}
	message := types.NewMessage(
		s.account.PublicKey,
		[]types.Instruction{
			instructions.UpdateConsulsInstruction(
				s.account.PublicKey, s.gravityContract, s.programID, currentConsuls, 3, uint64(roundId), solanaConsuls,
			),
		},
		res.Blockhash,
	)

	return message, nil
}
