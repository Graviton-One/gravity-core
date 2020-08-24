package deployer

import (
	"context"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/contracts"

	wavesHelper "github.com/Gravity-Tech/gravity-core/common/helpers/waves"
	wavesClient "github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/wavesplatform/gowaves/pkg/proto"
)

func DeployGravityWaves(client *wavesClient.Client, helper wavesHelper.ClientHelper, gravityScript []byte, consuls []string, bftValue int64, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) error {
	id, err := DeployWavesContract(client, gravityScript, chainId, secret, ctx)
	if err != nil {
		return err
	}

	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	id, err = DataWavesContract(client, chainId, secret, proto.DataEntries{
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
func DeployNebulaWaves(client *wavesClient.Client, helper wavesHelper.ClientHelper, nebulaScript []byte, gravityAddress string, subscriberAddress string,
	oracles []string, bftValue int64, dataType contracts.ExtractorType, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) error {

	id, err := DeployWavesContract(client, nebulaScript, chainId, secret, ctx)
	if err != nil {
		return err
	}

	err = <-helper.WaitTx(id, ctx)
	if err != nil {
		return err
	}

	id, err = DataWavesContract(client, chainId, secret, proto.DataEntries{
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

func DeploySubWaves(client *wavesClient.Client, helper wavesHelper.ClientHelper, subScript []byte, chainId byte, secret wavesCrypto.SecretKey, ctx context.Context) error {
	id, err := DeployWavesContract(client, subScript, chainId, secret, ctx)
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
