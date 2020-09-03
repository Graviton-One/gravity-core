package deployer

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployGravityEthereum(ethClient *ethclient.Client, consuls []common.Address, bftValue int64, ethPrivKey *ecdsa.PrivateKey) (common.Address, *contracts.Gravity, error) {
	transactOpt := bind.NewKeyedTransactor(ethPrivKey)
	address, _, gravity, err := contracts.DeployGravity(transactOpt, ethClient, consuls, big.NewInt(bftValue))
	if err != nil {
		return common.Address{}, nil, err
	}

	return address, gravity, nil
}

func DeployNebulaEthereum(ethClient *ethclient.Client, extractorType contracts.ExtractorType, gravityContract common.Address,
	subscriber common.Address, oracles []common.Address, bftValue int64, ethPrivKey *ecdsa.PrivateKey) (common.Address, *contracts.Nebula, error) {

	transactOpt := bind.NewKeyedTransactor(ethPrivKey)
	address, _, gravity, err := contracts.DeployNebula(transactOpt, ethClient, uint8(extractorType), gravityContract, oracles, big.NewInt(bftValue), subscriber)
	if err != nil {
		return common.Address{}, nil, err
	}

	return address, gravity, nil
}
