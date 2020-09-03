package deployer

import (
	"context"
	"time"

	wavesHelper "rebbit-hole/helpers"

	wavesClient "github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/wavesplatform/gowaves/pkg/proto"
)

func CreateToken(client *wavesClient.Client, helper wavesHelper.ClientHelper, name string, description string, quantity uint64, decimals byte, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) (string, error) {
	tx := &proto.IssueWithSig{
		Type:    proto.IssueTransaction,
		Version: 1,
		Issue: proto.Issue{
			Name:        name,
			Description: description,
			Fee:         100000000,
			Timestamp:   wavesClient.NewTimestampFromTime(time.Now()),
			Reissuable:  false,
			Quantity:    quantity,
			Decimals:    decimals,
			SenderPK:    wavesCrypto.GeneratePublicKey(secret),
		},
	}
	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}
	_, err = client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	err = <-helper.WaitTx(tx.ID.String(), ctx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func DeployLUPort(client *wavesClient.Client, helper wavesHelper.ClientHelper, luPortScript []byte, nebulaAddress string, assetId string, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) error {
	id, err := DeployWavesContract(client, luPortScript, chainId, secret, ctx)
	if err != nil {
		return err
	}
	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	id, err = DataWavesContract(client, chainId, secret, proto.DataEntries{
		&proto.StringDataEntry{
			Key:   "nebula_address",
			Value: nebulaAddress,
		},
		&proto.StringDataEntry{
			Key:   "asset_id",
			Value: assetId,
		},
	}, ctx)
	if err != nil {
		return err
	}

	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeployWavesContract(client *wavesClient.Client, contactScript []byte, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) (string, error) {
	tx := &proto.SetScriptWithProofs{
		Type:      proto.SetScriptTransaction,
		Version:   1,
		SenderPK:  wavesCrypto.GeneratePublicKey(secret),
		ChainID:   chainId,
		Script:    contactScript,
		Fee:       10000000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	}
	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}

	_, err = client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func DataWavesContract(client *wavesClient.Client, chainId byte, secret wavesCrypto.SecretKey, dataEntries proto.DataEntries, ctx context.Context) (string, error) {
	tx := &proto.DataWithProofs{
		Type:      proto.DataTransaction,
		Version:   1,
		SenderPK:  wavesCrypto.GeneratePublicKey(secret),
		Entries:   dataEntries,
		Fee:       10000000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	}

	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}

	_, err = client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
