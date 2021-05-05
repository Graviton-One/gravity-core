package adaptors

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/abi/solana/instructions"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"github.com/gorilla/websocket"
	"github.com/mr-tron/base58/base58"
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
	multisigAccount solana_common.PublicKey
	client          *solana.Client
	ghClient        *gravity.Client
	recentBlockHash string
}

func NewSolanaAdaptor(privKey []byte, nodeUrl string, custom map[string]interface{}, ghClient *gravity.Client) (*SolanaAdapter, error) {

	account := types.AccountFromPrivateKeyBytes(privKey)
	solClient := solana.NewClient(nodeUrl)
	gravityContract, ok := custom["gravity_contract"].(string)
	if !ok {
		zap.L().Error("Cannot parse gravity contract")
	}
	multisigAccount, ok := custom["multisig_account"].(string)
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
		multisigAccount: solana_common.PublicKeyFromString(multisigAccount),
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

	u := url.URL{Scheme: "ws", Host: "testnet.solana.com", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	req := `
		{
			"jsonrpc": "2.0",
			"id": 1,
			"method": "signatureSubscribe",
			"params": [
			  "%s",
			  {
				"commitment": "finalized"
			  }
			]
		  }
		`
	unsubscribeRequest := "{\"jsonrpc\":\"2.0\", \"id\":1, \"method\":\"signatureUnsubscribe\", \"params\":[%d]}"
	done := make(chan struct{})
	defer close(done)
	go func() {

		subscription := 0
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			a := struct {
				Method       string `json:"method"`
				Result       int    `json:"result"`
				Subscription int    `json:"subscription"`
			}{}
			json.Unmarshal(message, &a)
			switch a.Method {
			case "":
				subscription = a.Subscription
			case "signatureNotification":
				c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(unsubscribeRequest, subscription)))
				done <- struct{}{}
				c.Close()
				return
			}

		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(req, id)))
	<-done
	return nil
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
	zap.L().Sugar().Debug("Send msg: ", base58.Encode(serializedMessage))
	for key, sig := range signs {
		nkey := solana_common.PublicKey{}
		copy(nkey[:], key[1:33])
		solsigs[nkey] = sig
		zap.L().Sugar().Debug("Lsig: ", nkey.ToBase58(), " -> ", base58.Encode(sig))
	}
	solsigs[s.account.PublicKey] = selfSig
	zap.L().Sugar().Debug("Self sig: ", s.account.PublicKey.ToBase58(), " -> ", base58.Encode(selfSig))
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
	s.updateRecentBlockHash()
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
	zap.L().Sugar().Debug("msg: ", base58.Encode(serializedMessage))
	zap.L().Sugar().Debug("sig: ", base58.Encode(sign))
	return sign, nil
}

func (s *SolanaAdapter) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	panic("not implemented") // TODO: Implement
}

func (s *SolanaAdapter) LastRound(ctx context.Context) (uint64, error) {
	bft, _ := s.GetCurrentBFT()

	r, err := s.client.GetAccountInfo(s.gravityContract.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 8,
			Offset: 2 + 32 + uint64(bft*32),
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}

	sval, ok := r.Data.([]interface{})[0].(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return 0, err
	}

	val, err := base64.StdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}
	round := binary.LittleEndian.Uint64(val)
	return round, nil
}

func (s *SolanaAdapter) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	return false, nil //default mock
	//panic("not implemented") // TODO: Implement
}

//Custom solana methods

func (s *SolanaAdapter) GetCurrentConsuls() ([]solana_common.PublicKey, error) {
	bft, _ := s.GetCurrentBFT()
	length := uint64(bft * 32)
	r, err := s.client.GetAccountInfo(s.gravityContract.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: length,
			Offset: 34,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return []solana_common.PublicKey{}, err
	}
	sval, ok := r.Data.([]interface{})[0].(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return []solana_common.PublicKey{}, err
	}
	val, err := base64.StdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return []solana_common.PublicKey{}, err
	}

	sconsuls := SortablePubkey{}
	for i := 0; i < int(bft*32); i += 32 { // BFT=3
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
			//consuls = append(consuls, types.Account{}.PublicKey)
			continue
		}
		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		consuls = append(consuls, pubKey)
	}
	sort.Sort(&consuls)
	solanaConsuls := consuls.ToPubKeys()

	// res, err := s.client.GetRecentBlockhash()
	// if err != nil {
	// 	zap.L().Sugar().Error(err.Error())
	// 	return types.Message{}, err
	// }

	message := types.NewMessage(
		s.account.PublicKey,
		[]types.Instruction{
			instructions.UpdateConsulsInstruction(
				s.account.PublicKey, s.gravityContract, s.programID, s.multisigAccount, currentConsuls, uint8(len(currentConsuls)), uint64(roundId), solanaConsuls,
			),
		},
		s.recentBlockHash,
	)

	return message, nil
}

func (s *SolanaAdapter) updateRecentBlockHash() {
	res, err := s.client.GetRecentBlockhash()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return
	}
	s.recentBlockHash = res.Blockhash
}

func (s *SolanaAdapter) GetCurrentBFT() (byte, error) {
	r, err := s.client.GetAccountInfo(s.gravityContract.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 1,
			Offset: 33,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}
	sval, ok := r.Data.([]interface{})[0].(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return 0, err
	}
	val, err := base64.StdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}

	return val[0], nil
}
