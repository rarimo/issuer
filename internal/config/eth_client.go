package config

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type EthClientConfig struct {
	EthClient            *ethclient.Client
	StateStorageContract *common.Address
	PrivateKey           *ecdsa.PrivateKey
}

type ethClientConfigRaw struct {
	RpcURL               string `fig:"rpc_url,required"` //nolint
	StateStorageContract string `fig:"state_storage_contract,required"`
	PrivateKey           string `fig:"private_key,required"`
}

func (c *config) EthClient() *EthClientConfig {
	return c.ethClient.Do(func() interface{} {
		configRaw := ethClientConfigRaw{}
		err := figure.
			Out(&configRaw).
			From(kv.MustGetStringMap(c.getter, "ethereum")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		ethClient, err := ethclient.Dial(configRaw.RpcURL)
		if err != nil {
			panic(errors.Wrap(err, "failed to create an ethereum client"))
		}

		if !common.IsHexAddress(configRaw.StateStorageContract) {
			panic(errors.New("failed to parse state storage contract address"))
		}
		stateStorageAddress := common.HexToAddress(configRaw.StateStorageContract)

		privateKey, err := crypto.HexToECDSA(configRaw.PrivateKey)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse private key"))
		}

		return &EthClientConfig{
			EthClient:            ethClient,
			StateStorageContract: &stateStorageAddress,
			PrivateKey:           privateKey,
		}
	}).(*EthClientConfig)
}
