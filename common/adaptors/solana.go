package adaptors

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/abi/solana/instructions"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"github.com/Gravity-Tech/gravity-core/rpc"
	"github.com/mr-tron/base58/base58"
	"github.com/near/borsh-go"
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
func sumBytes(b []byte) int {
	sum := 0
	for i := 0; i < len(b); i++ {
		sum += int(b[i])
	}
	return sum
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
	account           types.Account
	programID         solana_common.PublicKey
	gravityContract   solana_common.PublicKey
	nebulaProgram     solana_common.PublicKey
	multisigAccount   solana_common.PublicKey
	client            *solana.Client
	ghClient          *gravity.Client
	recentBlockHashes map[string]string
	Bft               uint8
	oracleInterval    uint64

	ibportProgramAccount  solana_common.PublicKey
	ibportDataAccount     solana_common.PublicKey
	tokenProgramAddress   solana_common.PublicKey
	ibPortPDA             solana_common.PublicKey
	IBPortPDAtokenAccount solana_common.PublicKey
}
type SolanaAdapterOption func(*SolanaAdapter) error

func SolanaAdapterWithGhClient(ghClient *gravity.Client) SolanaAdapterOption {
	return func(s *SolanaAdapter) error {
		s.ghClient = ghClient
		return nil
	}
}

func SolanaAdapterWithCustom(custom map[string]interface{}) SolanaAdapterOption {
	return func(s *SolanaAdapter) error {
		gravityContract, ok := custom["gravity_contract"].(string)
		if ok {
			s.gravityContract = solana_common.PublicKeyFromString(gravityContract)
		}
		multisigAccount, ok := custom["multisig_account"].(string)
		if ok {
			s.multisigAccount = solana_common.PublicKeyFromString(multisigAccount)
		}
		nebulaContract, ok := custom["nebula_program"].(string)
		if ok {
			s.nebulaProgram = solana_common.PublicKeyFromString(nebulaContract)
		}

		programID, ok := custom["program_id"].(string)
		if ok {
			s.programID = solana_common.PublicKeyFromString(programID)
		}

		ibportProgramAccount, ok := custom["ib_port_program"].(string)
		if ok {
			s.ibportProgramAccount = solana_common.PublicKeyFromString(ibportProgramAccount)
		}

		ibportDataAccount, ok := custom["ib_port_data"].(string)
		if ok {
			s.ibportDataAccount = solana_common.PublicKeyFromString(ibportDataAccount)
		}

		tokenProgramAddress, ok := custom["token_program"].(string)
		if ok {
			s.tokenProgramAddress = solana_common.PublicKeyFromString(tokenProgramAddress)
		}

		ibPortPDA, ok := custom["ib_port_pda"].(string)
		if ok {
			s.ibPortPDA = solana_common.PublicKeyFromString(ibPortPDA)
		}
		ibPortPDAtokenAccount, ok := custom["ib_port_pda_token_account"].(string)
		if ok {
			s.IBPortPDAtokenAccount = solana_common.PublicKeyFromString(ibPortPDAtokenAccount)
		}

		return nil
	}
}
func NewSolanaAdaptor(privKey []byte, nodeUrl string, opts ...SolanaAdapterOption) (*SolanaAdapter, error) {

	account := types.AccountFromPrivateKeyBytes(privKey)
	solClient := solana.NewClient(nodeUrl)

	adapter := SolanaAdapter{
		client:  solClient,
		account: account,
	}
	adapter.recentBlockHashes = make(map[string]string)

	for _, opt := range opts {
		err := opt(&adapter)
		if err != nil {
			return nil, err
		}
	}
	if sumBytes(adapter.gravityContract[:]) != 0 {
		bft, err := adapter.GetCurrentBFT(context.Background())
		if err != nil {
			return nil, err
		}
		adapter.Bft = bft
	}
	if sumBytes(adapter.gravityContract[:]) != 0 {
		bft, err := adapter.GetCurrentBFT(context.Background())
		if err != nil {
			return nil, err
		}
		adapter.Bft = bft
	}
	return &adapter, nil
}

func (s *SolanaAdapter) GetHeight(ctx context.Context) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetHeight", r)
		}
	}()
	info, err := s.client.GetEpochInfo(ctx, solana.CommitmentFinalized)
	if err != nil {
		return 0, err
	}
	return uint64(info.BlockHeight), nil
}

func (s *SolanaAdapter) WaitTx(id string, ctx context.Context) error {
	time.Sleep(time.Second * 40)

	// u := url.URL{Scheme: "ws", Host: "testnet.solana.com", Path: "/"}
	// log.Printf("connecting to %s", u.String())

	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// 	log.Fatal("dial:", err)
	// }
	// defer c.Close()
	// req := `
	// 	{
	// 		"jsonrpc": "2.0",
	// 		"id": 1,
	// 		"method": "signatureSubscribe",
	// 		"params": [
	// 		  "%s",
	// 		  {
	// 			"commitment": "finalized"
	// 		  }
	// 		]
	// 	  }
	// 	`
	// unsubscribeRequest := "{\"jsonrpc\":\"2.0\", \"id\":1, \"method\":\"signatureUnsubscribe\", \"params\":[%d]}"
	// done := make(chan struct{})
	// defer close(done)
	// go func() {

	// 	subscription := 0
	// 	for {
	// 		_, message, err := c.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			return
	// 		}
	// 		log.Printf("recv: %s", message)
	// 		a := struct {
	// 			Method       string `json:"method"`
	// 			Result       int    `json:"result"`
	// 			Subscription int    `json:"subscription"`
	// 		}{}
	// 		json.Unmarshal(message, &a)
	// 		switch a.Method {
	// 		case "":
	// 			subscription = a.Subscription
	// 		case "signatureNotification":
	// 			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(unsubscribeRequest, subscription)))
	// 			done <- struct{}{}
	// 			c.Close()
	// 			return
	// 		}

	// 	}
	// }()

	// err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(req, id)))
	// <-done
	return nil
}

func (s *SolanaAdapter) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.account.PrivateKey, msg), nil
}
func (s *SolanaAdapter) SignHash(nebulaId account.NebulaId, intervalId uint64, pulseId uint64, hash []byte) ([]byte, error) {
	s.updateRecentBlockHash(context.Background(), "oracle")
	var validators []account.OraclesPubKey
	s.oracleInterval = intervalId
	oraclesMap, err := s.ghClient.BftOraclesByNebula(account.Solana, nebulaId)
	if err != nil {
		zap.L().Sugar().Debugf("BFT error: %s , \n %s", err, zap.Stack("trace").String)
		return []byte{}, nil
	}

	oracles := SortablePubkey{}
	for k, v := range oraclesMap {
		oracle, err := account.StringToOraclePubKey(k, v)
		if err != nil {
			return []byte{}, err
		}
		pubKey := solana_common.PublicKeyFromBytes(oracle[1:33])
		oracles = append(oracles, pubKey)
		validators = append(validators, oracle)
	}

	sort.Sort(&oracles)
	solanaOracles := oracles.ToPubKeys()

	if len(oracles) == 0 {
		zap.L().Debug("Oracles map is empty")
		return []byte{}, fmt.Errorf("Oracles map is empty")
	}

	sender := intervalId % uint64(len(solanaOracles))
	senderPubKey := solanaOracles[sender]
	msg, err := s.createAddPulseMessage(nebulaId, validators, pulseId, hash, senderPubKey)
	if err != nil {
		return []byte{}, err
	}
	serializedMessage, err := msg.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return []byte{}, err
	}

	sig, err := s.Sign(serializedMessage)
	zap.L().Sugar().Debugf("msg: %s", base58.Encode(serializedMessage))
	zap.L().Sugar().Debugf("sig: %s", base58.Encode(sig))
	return sig, err
}

func (s *SolanaAdapter) PubKey() account.OraclesPubKey {
	var pubKey account.OraclesPubKey
	copy(pubKey[:], append([]byte{0}, s.account.PublicKey.Bytes()[0:32]...))
	return pubKey
}

func (s *SolanaAdapter) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ValueType", r)
		}
	}()
	// n, err := s.getNebulaContractState(ctx, s.nebulaContract.ToBase58())
	// if err != nil {
	// 	zap.L().Error(err.Error())
	// 	return 0, err
	// }
	return abi.ExtractorType(abi.BytesType), nil //HardCode
}

func (s *SolanaAdapter) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in AddPulse", r)
		}
	}()

	oracles := SortablePubkey{}
	for _, v := range validators {

		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		oracles = append(oracles, pubKey)
	}

	sort.Sort(&oracles)
	solanaOracles := oracles.ToPubKeys()

	senderIndex := s.oracleInterval % uint64(len(oracles))
	senderPubKey := solanaOracles[senderIndex]
	msg, err := s.createAddPulseMessage(nebulaId, validators, pulseId, hash, senderPubKey)
	if err != nil {
		return "", err
	}
	serializedMessage, err := msg.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}
	solsigs := make(map[solana_common.PublicKey]types.Signature)
	for _, validator := range validators {
		vpk := solana_common.PublicKeyFromBytes(validator[1:33])
		sign, err := s.ghClient.Result(account.Solana, nebulaId, int64(pulseId), validator)
		if err != nil {
			zap.L().Sugar().Error(err.Error())
			continue
		}
		solsigs[vpk] = sign
		zap.L().Sugar().Debug("L sig: ", vpk.ToBase58(), " -> ", base58.Encode(sign))
	}
	selfSig, err := s.Sign(serializedMessage)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}
	zap.L().Sugar().Debug("Send msg: ", base58.Encode(serializedMessage))
	//solsigs[s.account.PublicKey] = selfSig
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

	txSig, err := s.client.SendRawTransaction(ctx, rawTx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	log.Println("txHash:", txSig)
	return txSig, nil
}

func (s *SolanaAdapter) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SendValueToSubs", r)
		}
	}()
	nid := solana_common.PublicKeyFromBytes(nebulaId[:])
	nst, err := s.getNebulaContractState(ctx, nid.ToBase58())
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return err
	}
	ids := nst.SubscriptionsMap.K
	dtype := uint8(0)
	for _, id := range ids {
		zap.L().Sugar().Debug("IDs iterate", id)
		val := []byte{}
		switch value.Type {
		case extractor.Int64:
			zap.L().Sugar().Debugf("SendIntValueToSubs")
			v, err := strconv.ParseInt(value.Value, 10, 64)
			if err != nil {
				return err
			}
			val = make([]byte, 8)
			binary.LittleEndian.PutUint64(val, uint64(v))
			dtype = 0
		case extractor.String:
			zap.L().Sugar().Debugf("SendStringValueToSubs")
			val = []byte(value.Value)
			dtype = 1
		case extractor.Base64:
			//println(value.Value)
			v, err := base64.StdEncoding.DecodeString(value.Value)
			if err != nil {
				return err
			}
			val = v
			dtype = 2
		}
		if len(val) > 0 {
			msg, err := s.createSendValueToSubsMessage(nebulaId, pulseId-1, dtype, val, id)
			if err != nil {
				return err
			}
			serializedMessage, err := msg.Serialize()
			if err != nil {
				zap.L().Sugar().Error(err.Error())
				return err
			}
			solsigs := make(map[solana_common.PublicKey]types.Signature)
			selfSig, err := s.Sign(serializedMessage)
			if err != nil {
				zap.L().Sugar().Error(err.Error())
				return err
			}
			zap.L().Sugar().Debug("Send msg: ", base58.Encode(serializedMessage))
			solsigs[s.account.PublicKey] = selfSig
			zap.L().Sugar().Debug("Self sig: ", s.account.PublicKey.ToBase58(), " -> ", base58.Encode(selfSig))
			tx, err := types.CreateTransaction(msg, solsigs)
			if err != nil {
				return err
			}
			rawTx, err := tx.Serialize()
			if err != nil {
				zap.L().Sugar().Error(err.Error())
				return err
			}
			zap.L().Sugar().Debug("SendValueToSubs(Base64): ", base64.StdEncoding.EncodeToString(rawTx))
			txSig, err := s.client.SendRawTransaction(ctx, rawTx)
			if err != nil {
				zap.L().Sugar().Error(err.Error())
				return err
			}

			log.Println("txHash:", txSig)
		}
	}
	return nil
}

func (s *SolanaAdapter) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SetOraclesToNebula", r)
		}
	}()

	zap.L().Sugar().Debugf("SetOraclesToNebula [%s]", solana_common.PublicKeyFromBytes(nebulaId[:]).ToBase58())
	customParams, err := rpc.GlobalClient.NebulaCustomParams(nebulaId, account.Solana)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	nebulaContract_interface, ok := customParams["nebula_contract"]
	if !ok {
		return "", fmt.Errorf("Data account for nebula not declared")
	}
	nebulaContract := nebulaContract_interface.(string)

	n, err := s.getNebulaContractState(ctx, nebulaContract)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}
	msg, err := s.createUpdateOraclesMessage(ctx, nebulaId, oracles, round, n.Bft, customParams, s.account.PublicKey)
	if err != nil {
		zap.L().Error(err.Error())
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
	//solsigs[s.account.PublicKey] = selfSig
	zap.L().Sugar().Debug("Self sig: ", s.account.PublicKey.ToBase58(), " -> ", base58.Encode(selfSig))
	tx, err := types.CreateTransaction(msg, solsigs)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}
	rawTx, err := tx.Serialize()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	txSig, err := s.client.SendRawTransaction(ctx, rawTx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	log.Println("txHash:", txSig)
	return txSig, nil
}

func (s *SolanaAdapter) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SendConsulsToGravityCOntract", r)
		}
	}()
	msg, err := s.createUpdateConsulsMessage(ctx, newConsulsAddresses, round, s.account.PublicKey)
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

	txSig, err := s.client.SendRawTransaction(ctx, rawTx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return "", err
	}

	log.Println("txHash:", txSig)
	return txSig, nil
}

func (s *SolanaAdapter) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64, sender account.OraclesPubKey) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SignConsuls", r)
		}
	}()
	s.updateRecentBlockHash(context.Background(), "update_consuls")
	senderPubKey := solana_common.PublicKeyFromBytes(sender[1:33])
	zap.L().Sugar().Debugf("Sender pubkey: %s", senderPubKey.ToBase58())
	msg, err := s.createUpdateConsulsMessage(context.Background(), consulsAddresses, roundId, senderPubKey)
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

func (s *SolanaAdapter) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, round int64, sender account.OraclesPubKey) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in SignOracles", r)
		}
	}()
	s.updateRecentBlockHash(context.Background(), "update_oracles")
	senderPubKey := solana_common.PublicKeyFromBytes(sender[1:33])
	customParams, err := rpc.GlobalClient.NebulaCustomParams(nebulaId, account.Solana)
	if err != nil {
		return []byte{}, err
	}
	zap.L().Sugar().Debug("Custom params ", customParams)
	// nebulaProgram, ok := customParams["nebula_program"].(string)
	// if !ok {
	// 	return []byte{}, fmt.Errorf("Nebula account for nebula not declared")
	// }
	dataAccount := solana_common.PublicKeyFromBytes(nebulaId[:])
	n, err := s.getNebulaContractState(context.Background(), dataAccount.ToBase58())
	if err != nil {
		return nil, err
	}
	new_oracles := oracles
	if len(new_oracles) == 0 {
		old_oracles := n.Oracles
		for _, or := range old_oracles {
			new_or := account.OraclesPubKey{}
			copy(new_or[:], append([]byte{0}, or[0:32]...))
			new_oracles = append(new_oracles, &new_or)
		}
	}
	msg, err := s.createUpdateOraclesMessage(context.Background(), nebulaId, new_oracles, round, n.Bft, customParams, senderPubKey)
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

func (s *SolanaAdapter) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in LatPulseId", r)
		}
	}()
	nid := solana_common.PublicKeyFromBytes(nebulaId[:])
	n, err := s.getNebulaContractState(ctx, nid.ToBase58())
	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}
	return n.LastPulseId, nil
}

func (s *SolanaAdapter) LastRound(ctx context.Context) (uint64, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in LastRound", r)
		}
	}()
	gs, err := s.getGravityContractState(ctx)

	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}

	return gs.LastRound, nil
}

func (s *SolanaAdapter) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	return false, nil //default mock
}

//Custom solana methods

func (s *SolanaAdapter) GetCurrentConsuls(ctx context.Context) ([]solana_common.PublicKey, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetCurrentConsuls", r)
		}
	}()
	gs, err := s.getGravityContractState(ctx)

	if err != nil {
		zap.L().Error(err.Error())
		return []solana_common.PublicKey{}, err
	}

	sconsuls := SortablePubkey{}
	for _, c := range gs.Consuls {
		sconsuls = append(sconsuls, c)
	}
	sort.Sort(&sconsuls)
	return sconsuls, nil
}

func (s *SolanaAdapter) createUpdateConsulsMessage(ctx context.Context, consulsAddresses []*account.OraclesPubKey, roundId int64, sender solana_common.PublicKey) (types.Message, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in createUpdateConsulsMesssage", r)
		}
	}()
	currentConsuls, err := s.GetCurrentConsuls(ctx)
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

	message := types.NewMessage(
		sender,
		[]types.Instruction{
			instructions.UpdateConsulsInstruction(
				sender, s.gravityContract, s.programID, s.multisigAccount, currentConsuls, uint8(len(currentConsuls)), uint64(roundId), solanaConsuls,
			),
		},
		s.recentBlockHashes["update_consuls"],
	)

	return message, nil
}

func (s *SolanaAdapter) createUpdateOraclesMessage(ctx context.Context, nebulaId account.NebulaId, oraclesAddresses []*account.OraclesPubKey, roundId int64, Bft uint8, customParams storage.NebulaCustomParams, sender solana_common.PublicKey) (types.Message, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in createUpdateOraclesMessage", r)
		}
	}()
	currentConsuls, err := s.GetCurrentConsuls(ctx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return types.Message{}, err
	}

	consuls := SortablePubkey{}
	for _, v := range currentConsuls {
		consuls = append(consuls, v)
	}
	sort.Sort(&consuls)
	solanaConsuls := consuls.ToPubKeys()

	newOracles := SortablePubkey{}
	for _, v := range oraclesAddresses {
		if v == nil {
			continue
		}
		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		newOracles = append(newOracles, pubKey)
	}
	sort.Sort(&newOracles)
	solanaOracles := newOracles.ToPubKeys()
	nebulaDataAccount := solana_common.PublicKeyFromBytes(nebulaId[:])
	// customParams, err := rpc.GlobalClient.NebulaCustomParams(nebulaId, account.Solana)
	// if err != nil {
	// 	return types.Message{}, err
	// }

	multisigAccount_interface, ok := customParams["multisig_account"]
	if !ok {
		return types.Message{}, fmt.Errorf("Multisig account for nebula data account [%s] not declared", nebulaDataAccount.ToBase58())
	}
	multisigAccount := multisigAccount_interface.(string)

	nebulaProgram_interface, ok := customParams["nebula_program"]
	if !ok {
		return types.Message{}, fmt.Errorf("Data account for nebula data account [%s] not declared", nebulaDataAccount.ToBase58())
	}
	nebulaProgram := nebulaProgram_interface.(string)

	message := types.NewMessage(
		sender,
		[]types.Instruction{
			instructions.NebulaUpdateOraclesInstruction(
				sender,
				solana_common.PublicKeyFromString(nebulaProgram),
				nebulaDataAccount,
				solana_common.PublicKeyFromString(multisigAccount),
				solanaConsuls,
				uint64(roundId),
				solanaOracles,
				Bft,
			),
		},
		s.recentBlockHashes["update_oracles"],
	)
	zap.L().Sugar().Debug("Message created: ", message)
	return message, nil
}

func (s *SolanaAdapter) createAddPulseMessage(nebulaId account.NebulaId, validators []account.OraclesPubKey, pulseId uint64, hash []byte, sender solana_common.PublicKey) (types.Message, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in createAddPulseMessage", r)
		}
	}()
	vals := SortablePubkey{}
	for _, v := range validators {
		pubKey := solana_common.PublicKeyFromBytes(v[1:33])
		vals = append(vals, pubKey)
	}
	nebulaDataAccount := solana_common.PublicKeyFromBytes(nebulaId[:])
	sort.Sort(&vals)
	solanaValidators := vals.ToPubKeys()
	message := types.NewMessage(
		sender,
		[]types.Instruction{
			instructions.NebulaAddPulseInstruction(
				sender, s.nebulaProgram, s.nebulaProgram, s.multisigAccount, nebulaDataAccount, solanaValidators, pulseId, hash,
			),
		},
		s.recentBlockHashes["oracle"],
	)
	zap.L().Sugar().Debugf("Block hash %s", s.recentBlockHashes["oracle"])
	return message, nil
}

func RecipientFromByteArray(byteArray []byte) []byte {
	// 'm' (1 byte) + swapId (16 bytes) + float (8 bytes)
	offset := 1 + 16 + 8
	return byteArray[offset : offset+32]
}

func (s *SolanaAdapter) createSendValueToSubsMessage(nebulaId account.NebulaId, pulseId uint64, DataType uint8, value []byte, id [16]byte) (types.Message, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in createSendValueToSubsMessage", r)
		}
	}()
	recentBlockHash, err := s.client.GetRecentBlockhash(context.Background())
	if err != nil {
		return types.Message{}, err
	}
	nebulaDataAccount := solana_common.PublicKeyFromBytes(nebulaId[:])
	recipient := solana_common.PublicKeyFromBytes(RecipientFromByteArray(value))

	resp, err := s.client.GetAccountInfo(context.Background(), recipient.ToBase58(), solana.GetAccountInfoConfig{
		Encoding: "base64",
	})
	if err != nil {
		return types.Message{}, err
	}
	recipientOwner := solana_common.PublicKeyFromString(resp.Owner)
	message := types.NewMessage(
		s.account.PublicKey,
		[]types.Instruction{
			instructions.NebulaSendValueToSubsInstruction(
				s.account.PublicKey, s.nebulaProgram, s.nebulaProgram,
				nebulaDataAccount, s.multisigAccount,
				s.ibportProgramAccount, s.ibportDataAccount,
				s.tokenProgramAddress, recipient, s.ibPortPDA, recipientOwner,
				s.IBPortPDAtokenAccount, DataType, value, pulseId, id,
			),
		},
		recentBlockHash.Blockhash,
	)
	return message, nil
}

func (s *SolanaAdapter) updateRecentBlockHash(ctx context.Context, key string) {
	res, err := s.client.GetRecentBlockhash(ctx)
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return
	}
	s.recentBlockHashes[key] = res.Blockhash
	zap.L().Sugar().Debugf("New ReccentBlockHash: %s", s.recentBlockHashes[key])
}

// func (s *SolanaAdapter) updateOraclesRecentBlockHash(ctx context.Context) {
// 	res, err := s.client.GetRecentBlockhash(ctx)
// 	if err != nil {
// 		zap.L().Sugar().Error(err.Error())
// 		return
// 	}
// 	s.oraclesRecentBlockHash = res.Blockhash
// 	zap.L().Sugar().Debugf("New Oracles ReccentBlockHash: %s", s.recentBlockHash)
// }

func (s *SolanaAdapter) GetCurrentBFT(ctx context.Context) (byte, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetCurrentBFT", r)
		}
	}()
	gs, err := s.getGravityContractState(ctx)

	if err != nil {
		zap.L().Error(err.Error())
		return 0, err
	}

	return gs.Bft, nil
}

func (s *SolanaAdapter) getCurrentOracles() ([]solana_common.PublicKey, error) {
	return []solana_common.PublicKey{}, nil
}

func (s *SolanaAdapter) getNebulaContractState(ctx context.Context, stateAccount string) (*instructions.NebulaContract, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getNebulaContractState", r)
		}
	}()
	//nid := base58.Encode(nebulaId[:])
	zap.L().Sugar().Debugf("gettingNebulaState: %s", stateAccount)
	r, err := s.client.GetAccountInfo(ctx, stateAccount, solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 2000,
			Offset: 0,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if r.Data == nil {
		zap.L().Error("Invalid account data")
		return nil, fmt.Errorf("empty account data")
	}
	sval, ok := r.Data.([]interface{})[0].(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return nil, err
	}

	val, err := base64.StdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	n := instructions.NewNebulaContract()
	err = borsh.Deserialize(&n, val)

	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &n, nil
}

func (s *SolanaAdapter) getGravityContractState(ctx context.Context) (*instructions.GravityContract, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in getNebulaContractState", r)
		}
	}()
	gid := s.gravityContract.ToBase58()
	r, err := s.client.GetAccountInfo(ctx, gid, solana.GetAccountInfoConfig{
		Encoding: "base64",
		DataSlice: solana.GetAccountInfoConfigDataSlice{
			Length: 299,
			Offset: 0,
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if r.Data == nil {
		zap.L().Error("Invalid account data")
		return nil, fmt.Errorf("empty account data")
	}
	sval, ok := r.Data.([]interface{})[0].(string)
	if !ok {
		zap.L().Error("Invalid account data")
		return nil, err
	}

	val, err := base64.StdEncoding.DecodeString(sval)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	n := instructions.NewGravityContract()
	err = borsh.Deserialize(&n, val)

	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &n, nil
}
