package publisher

import (
	"bytes"
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

func (p *publisher) waitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, p.ethClient, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get mined tx")
	}

	if receipt.Status == types.ReceiptStatusFailed {
		reason, err := checkTxErrorReason(ctx, tx, p.ethClient, p.address)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get tx error reason")
		}
		return nil, errors.Wrap(ErrTransactionFailed, reason)
	}

	return receipt, nil
}

func (p *publisher) retryChainCall(
	ctx context.Context,
	contractCall func(signer *bind.TransactOpts) error,
) error {
	signer, err := p.newSigner(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get new signer")
	}

	for {
		err := contractCall(signer)
		if err == nil {
			break
		}
		if !errors.Is(err, core.ErrNonceTooLow) {
			return errors.Wrap(err, "failed to call contract")
		}

		signer.Nonce.Add(signer.Nonce, big.NewInt(1))
	}
	return nil
}

func (p *publisher) newSigner(ctx context.Context) (*bind.TransactOpts, error) {
	nonce, err := p.ethClient.PendingNonceAt(ctx, p.address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get nonce")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(p.privateKey, p.chainID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction signer")
	}
	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}

func checkTxErrorReason(
	ctx context.Context,
	tx *types.Transaction,
	ethClient *ethclient.Client,
	addr common.Address,
) (string, error) {
	msg := ethereum.CallMsg{
		From:     addr,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to CallContract")
	}

	return unpackError(res)
}

func unpackError(result []byte) (string, error) {
	if len(result) < 4 || !bytes.Equal(result[:4], invalidSignatureError) {
		return "tx result is not Error(string)", errors.New("TX result not of type Error(string)")
	}

	values, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "invalid tx result", errors.Wrap(err, "unpacking revert reason")
	}

	return values[0].(string), nil
}
