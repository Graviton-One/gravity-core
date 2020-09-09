package deployer

import (
	"context"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func DeployGravityWaves(wClient *client.Client, helper helpers.ClientHelper, gravityScript []byte, consuls []string, bftValue int64, chainId byte, secret crypto.SecretKey, ctx context.Context) error {
	id, err := DeployWavesContract(wClient, gravityScript, chainId, secret, ctx)
	if err != nil {
		return err
	}

	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	id, err = DataWavesContract(wClient, chainId, secret, proto.DataEntries{
		&proto.StringDataEntry{
			Key:   "consuls",
			Value: strings.Join(consuls, ","),
		},
		&proto.IntegerDataEntry{
			Key:   "bft_coefficient",
			Value: bftValue,
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
func DeployNebulaWaves(wClient *client.Client, helper helpers.ClientHelper, nebulaScript []byte, gravityAddress string, subscriberAddress string,
	oracles []string, bftValue int64, dataType contracts.ExtractorType, chainId byte, secret crypto.SecretKey, ctx context.Context) error {

	id, err := DeployWavesContract(wClient, nebulaScript, chainId, secret, ctx)
	if err != nil {
		return err
	}

	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	id, err = DataWavesContract(wClient, chainId, secret, proto.DataEntries{
		&proto.StringDataEntry{
			Key:   "oracles",
			Value: strings.Join(oracles, ","),
		},
		&proto.IntegerDataEntry{
			Key:   "bft_coefficient",
			Value: bftValue,
		}, &proto.StringDataEntry{
			Key:   "subscriber_address",
			Value: subscriberAddress,
		}, &proto.StringDataEntry{
			Key:   "gravity_contract",
			Value: gravityAddress,
		},
		&proto.IntegerDataEntry{
			Key:   "type",
			Value: int64(dataType),
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

func DeploySubWaves(wClient *client.Client, helper helpers.ClientHelper, subScript []byte, chainId byte, secret crypto.SecretKey, ctx context.Context) error {
	id, err := DeployWavesContract(wClient, subScript, chainId, secret, ctx)
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
