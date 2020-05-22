package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"gh-node/api/gravity"
	"gh-node/config"
	"gh-node/extractors"
	"gh-node/helpers"
	"gh-node/keys"
	"gh-node/transaction"
	"strconv"
	"strings"
	"time"

	"github.com/mr-tron/base58"

	"github.com/wavesplatform/gowaves/pkg/client"

	"golang.org/x/net/context"

	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	proto "github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	DefaultConfigFileName = "config.json"
	Rounds                = 4
)

func logErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Error: %s\n", err.Error())
}
func main() {
	var confFileName, seedString string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.StringVar(&seedString, "seed", "", "set seed")
	flag.Parse()

	ctx := context.Background()
	wCrypto := wavesplatform.NewWavesCrypto()
	seed := wavesplatform.Seed(seedString)
	secret, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
	if err != nil {
		panic(err)
	}
	pubKey := crypto.GeneratePublicKey(secret)

	cfg, err := config.Load(confFileName)
	if err != nil {
		panic(err)
	}

	ghClient := gravity.NewClient(cfg.GHNodeURL)
	nebulaId, err := hex.DecodeString(cfg.NebulaId)
	if err != nil {
		panic(err)
	}

	validatorKey := keys.FormValidatorKey(nebulaId, pubKey.Bytes())
	_, err = ghClient.GetKey(validatorKey, ctx)
	if err == gravity.KeyNotFound {
		tx, err := transaction.New(pubKey.Bytes(), transaction.AddValidator, secret, append(nebulaId, pubKey.Bytes()...))
		if err != nil {
			panic(err)
		}

		err = ghClient.SendTx(tx, ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Add validator tx id: %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	validatorPrefix := strings.Join([]string{string(keys.ValidatorKey), hex.EncodeToString(nebulaId)}, "_")
	values, err := ghClient.GetByPrefix(validatorPrefix, ctx)
	if err != nil {
		panic(err)
	}
	var validators [][]byte
	var myRound = 0
	for k, _ := range values {
		keyParts := strings.Split(k, "_")
		validator, err := hex.DecodeString(keyParts[2])
		if err != nil {
			continue
		}

		if bytes.Compare(validator, pubKey.Bytes()) == 0 {
			myRound = len(validators)
		}

		validators = append(validators, validator)
	}

	wavesClient, err := client.NewClient(client.Options{ApiKey: "", BaseUrl: cfg.WavesNodeUrl})
	if err != nil {
		panic(err)
	}

	var lastWavesHeight uint64
	var lastGHHeight int64

	var extractor extractors.PriceExtractor = &extractors.BinanceExtractor{}
	var lastCommitPrice string
	var commitHeight int64
	var commitHash []byte
	var resultHash []byte
	for {
		price, err := extractor.PriceNow()
		logErr(err)
		wavesHeight, _, err := wavesClient.Blocks.Height(ctx)
		logErr(err)
		if lastWavesHeight != wavesHeight.Height {
			fmt.Printf("Waves Height: %d\n", wavesHeight.Height)
			lastWavesHeight = wavesHeight.Height
		}

		block, err := ghClient.GetBlock(ctx)
		logErr(err)
		ghHeight, err := strconv.ParseInt(block.Result.Block.Header.Height, 10, 64)
		logErr(err)

		if lastGHHeight != ghHeight {
			fmt.Printf("GH Height: %d\n", ghHeight)
			lastGHHeight = ghHeight
		}

		if ghHeight%Rounds == 0 {
			commitKey := keys.FormCommitKey(nebulaId, wavesHeight.Height, pubKey.Bytes())
			_, err = ghClient.GetKey(commitKey, ctx)
			if err == gravity.KeyNotFound {
				lastCommitPrice = fmt.Sprintf("%.2f", price)
				commit := sha256.Sum256([]byte(lastCommitPrice))

				fmt.Printf("Commit: %s - %s \n", lastCommitPrice, hex.EncodeToString(commit[:]))
				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, wavesHeight.Height)

				tx, err := transaction.New(pubKey.Bytes(), transaction.Commit, secret, append(nebulaId, append(heightBytes, commit[:]...)...))
				if err != nil {
					panic(err)
				}

				err = ghClient.SendTx(tx, ctx)
				logErr(err)
				fmt.Printf("Commit txId: %s\n", tx.Id)

				commitHeight = ghHeight
				commitHash = commit[:]
			} else {
				logErr(err)
			}
		}

		if commitHeight != 0 && ghHeight%Rounds == 1 {
			revealKey := keys.FormRevealKey(nebulaId, wavesHeight.Height, commitHash)

			_, err = ghClient.GetKey(revealKey, ctx)
			if err == gravity.KeyNotFound {
				fmt.Printf("Reveal: %s - %s \n", lastCommitPrice, hex.EncodeToString(commitHash[:]))
				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, wavesHeight.Height)
				var args []byte
				args = append(args, commitHash[:]...)
				args = append(args, nebulaId...)
				args = append(args, heightBytes...)
				args = append(args, lastCommitPrice...)

				tx, err := transaction.New(pubKey.Bytes(), transaction.Reveal, secret, args)
				if err != nil {
					panic(err)
				}

				err = ghClient.SendTx(tx, ctx)
				logErr(err)
				fmt.Printf("Reveal txId: %s\n", tx.Id)
				commitHeight = 0
			} else {
				logErr(err)
			}
		}

		if ghHeight%Rounds == 2 {
			signKey := keys.FormSignResultKey(nebulaId, wavesHeight.Height, pubKey.Bytes())

			_, err = ghClient.GetKey(signKey, ctx)
			if err == gravity.KeyNotFound {
				prefix := strings.Join([]string{string(keys.RevealKey), hex.EncodeToString(nebulaId), fmt.Sprintf("%d", wavesHeight.Height)}, "_")

				values, err := ghClient.GetByPrefix(prefix, ctx)
				if err != nil {
					panic(err)
				}

				var reveals []float64
				for _, v := range values {
					value, err := strconv.ParseFloat(string(v), 64)
					if err != nil {
						continue
					}
					reveals = append(reveals, value)
				}
				var average float64
				for _, v := range reveals {
					average += v
				}
				average = average / float64(len(reveals))
				result := fmt.Sprintf("%.2f", average)
				resultHashByte32 := sha256.Sum256([]byte(result))
				resultHash = resultHashByte32[:]
				fmt.Printf("Result hash: %s \n", hex.EncodeToString(resultHash))
				signBytes, err := crypto.Sign(secret, resultHash)
				logErr(err)

				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, wavesHeight.Height)
				var args []byte
				args = append(args, nebulaId[:]...)
				args = append(args, heightBytes...)
				args = append(args, resultHash...)
				args = append(args, signBytes.Bytes()...)

				tx, err := transaction.New(pubKey.Bytes(), transaction.SignResult, secret, args)
				if err != nil {
					panic(err)
				}

				err = ghClient.SendTx(tx, ctx)
				logErr(err)
				fmt.Printf("Sign result txId: %s\n", tx.Id)
			}
		}
		if resultHash != nil && ghHeight%Rounds == 3 && int(wavesHeight.Height)%len(validators) == myRound {
			helperWaves := helpers.New(cfg.WavesNodeUrl, "")
			state, err := helperWaves.GetStateByAddressAndKey(cfg.NebulaContract, fmt.Sprintf("%d", wavesHeight.Height))
			logErr(err)
			if state == nil {
				funcArgs := new(proto.Arguments)
				funcArgs.Append(proto.StringArgument{
					Value: base58.Encode(resultHash),
				})
				bft := int(float32(len(validators)) * 0.7)
				realSignCount := 0
				var signs []string
				for _, validator := range validators {
					sign, err := ghClient.GetKey(keys.FormSignResultKey(nebulaId, wavesHeight.Height, validator), ctx)
					if err != nil {
						signs = append(signs, "nil")
						continue
					}
					signs = append(signs, base58.Encode(sign))
					realSignCount++
				}
				funcArgs.Append(proto.StringArgument{
					Value: strings.Join(signs, ","),
				})

				if realSignCount >= bft {
					asset, err := proto.NewOptionalAssetFromString("WAVES")
					logErr(err)
					contract, err := proto.NewRecipientFromString(cfg.NebulaContract)
					logErr(err)
					tx := &proto.InvokeScriptWithProofs{
						Type:            proto.InvokeScriptTransaction,
						Version:         1,
						SenderPK:        pubKey,
						ChainID:         'T',
						ScriptRecipient: contract,
						FunctionCall: proto.FunctionCall{
							Name:      "write",
							Arguments: *funcArgs,
						},
						Payments:  nil,
						FeeAsset:  *asset,
						Fee:       500000,
						Timestamp: client.NewTimestampFromTime(time.Now()),
					}

					err = tx.Sign('T', secret)
					logErr(err)

					_, err = wavesClient.Transactions.Broadcast(ctx, tx)
					logErr(err)

					fmt.Printf("Tx finilize: %s \n", tx.ID)
				}
			}
		}
		time.Sleep(time.Duration(cfg.Timeout) * time.Second)
	}
}
