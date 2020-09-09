package deployer

import (
	"context"
	"time"

	"rebbit-hole/helpers"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/wavesplatform/gowaves/pkg/proto"
)

func CreateToken(wClient *client.Client, helper helpers.ClientHelper, name string, description string, quantity uint64, decimals byte, chainId byte, secret crypto.SecretKey, ctx context.Context) (string, error) {
	tx := &proto.IssueWithSig{
		Type:    proto.IssueTransaction,
		Version: 1,
		Issue: proto.Issue{
			Name:        name,
			Description: description,
			Fee:         100000000,
			Timestamp:   client.NewTimestampFromTime(time.Now()),
			Reissuable:  false,
			Quantity:    quantity,
			Decimals:    decimals,
			SenderPK:    crypto.GeneratePublicKey(secret),
		},
	}
	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}
	_, err = wClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	err = <-helper.WaitTx(tx.ID.String(), ctx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func DeployLUPort(client *client.Client, helper helpers.ClientHelper, luPortScript []byte, nebulaAddress string, assetId string, chainId byte, secret crypto.SecretKey, ctx context.Context) error {
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

func DeployWavesContract(wClient *client.Client, contactScript []byte, chainId byte, secret crypto.SecretKey, ctx context.Context) (string, error) {
	tx := &proto.SetScriptWithProofs{
		Type:      proto.SetScriptTransaction,
		Version:   1,
		SenderPK:  crypto.GeneratePublicKey(secret),
		ChainID:   chainId,
		Script:    contactScript,
		Fee:       10000000,
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}
	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}

	_, err = wClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func DataWavesContract(wClient *client.Client, chainId byte, secret crypto.SecretKey, dataEntries proto.DataEntries, ctx context.Context) (string, error) {
	tx := &proto.DataWithProofs{
		Type:      proto.DataTransaction,
		Version:   1,
		SenderPK:  crypto.GeneratePublicKey(secret),
		Entries:   dataEntries,
		Fee:       10000000,
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}

	err := tx.Sign(chainId, secret)
	if err != nil {
		return "", err
	}

	_, err = wClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
