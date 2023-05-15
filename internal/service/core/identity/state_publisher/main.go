package statepublisher

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/pkg/errors"

	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state_publisher/contracts"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/zkp"
)

type Publisher interface {
	Run(ctx context.Context)
}

func New(cfg *Config, state *state.IdentityState) (Publisher, error) {
	stateStorageContract, err := contracts.NewStateStorage(*cfg.EthConfig.StateStorageContract, cfg.EthConfig.EthClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new state store contract")
	}

	chainID, err := cfg.EthConfig.EthClient.ChainID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chainID")
	}

	publisherInstance := &publisher{
		log:                cfg.Log,
		stateStoreContract: stateStorageContract,
		state:              state,
		publishPeriod:      cfg.StatePublisher.PublishPeriod,
		retryPeriod:        cfg.StatePublisher.RetryPeriod,

		ethClient:  cfg.EthConfig.EthClient,
		privateKey: cfg.EthConfig.PrivateKey,
		address:    crypto.PubkeyToAddress(cfg.EthConfig.PrivateKey.PublicKey),
		chainID:    chainID,
	}

	err = publisherInstance.processPreviousSessionCommit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact pending queue")
	}

	return publisherInstance, nil
}

func (p *publisher) sendTransaction(
	ctx context.Context,
	tsInfo *state.StateTransitionInfo,
	committedState *data.CommittedState,
) (*types.Transaction, error) {
	zkpArgs, err := parseZKPArgs(tsInfo.ZKProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse SKProof to contract readable format")
	}

	var tx *types.Transaction
	err = p.retryChainCall(ctx, func(signer *bind.TransactOpts) error {
		tx, err = p.stateStoreContract.TransitState(
			signer,
			tsInfo.Identifier.BigInt(),
			tsInfo.LatestState.BigInt(),
			tsInfo.NewState.BigInt(),
			tsInfo.IsOldStateGenesis,
			zkpArgs.proofA, zkpArgs.proofB, zkpArgs.proofC,
		)
		if err != nil {
			if err := p.setStatusFailed("", err.Error(), committedState); err != nil {
				return errors.Wrap(err, "failed to set status failed in db")
			}

			return errors.Wrap(err, "failed to call transit state contract method")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to call retryChainCall")
	}

	committedState.Status = data.StatusProcessing
	err = p.state.DB.CommittedStatesQ().Update(committedState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update committed state status to processing")
	}

	return tx, nil
}

func (p *publisher) waitTransaction(
	ctx context.Context,
	committedState *data.CommittedState,
	tx *types.Transaction,
) error {
	receipt, err := p.waitMined(ctx, tx)
	if err != nil {
		if errors.Is(err, ErrTransactionFailed) {
			p.log.WithField("reason", err.Error()).Info("Failed to wait mined tx")
			err = p.setStatusFailed(tx.Hash().Hex(), err.Error(), committedState)
			if err != nil {
				return errors.Wrap(err, "failed to set status failed in db")
			}
			return nil
		}

		err = p.setStatusFailed(tx.Hash().Hex(), err.Error(), committedState)
		if err != nil {
			return errors.Wrap(err, "failed to set status failed in db")
		}

		return errors.Wrap(err, "failed to wait mined tx")
	}

	block, err := p.ethClient.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		p.log.WithError(err).Error("Failed to get ethereum block by number")
	}

	if block != nil {
		committedState.TxID = receipt.TxHash.Hex()
		committedState.BlockTimestamp = block.Time()
		committedState.BlockNumber = block.NumberU64()
	}

	err = p.setStatusCompleted(ctx, receipt, committedState)
	if err != nil {
		return errors.Wrap(err, "failed to set status completed")
	}

	p.log.WithField("tx_hash", tx.Hash().Hex()).Info("State was successfully transited")
	return nil
}

// nolint
func parseZKPArgs(zkp *zkp.ZKProof) (*contractReadableZKP, error) {
	a, b, c, err := zkp.ProofToBigInts()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ZKProof to big ints")
	}
	proofA := [2]*big.Int{a[0], a[1]}
	proofB := [2][2]*big.Int{
		{b[0][1], b[0][0]},
		{b[1][1], b[1][0]},
	}
	proofC := [2]*big.Int{c[0], c[1]}

	return &contractReadableZKP{
		proofA: proofA,
		proofB: proofB,
		proofC: proofC,
	}, nil
}

func (p *publisher) processPreviousSessionCommit() error {
	unprocessedStateList, err := p.state.DB.CommittedStatesQ().WhereStatus(data.StatusProcessing).Select()
	if err != nil {
		return errors.Wrap(err, "failed to select unprocessed committed states from db")
	}

	for _, unprocessedState := range unprocessedStateList {
		unprocessedStateHash, err := getStateHash(unprocessedState)
		if err != nil {
			return errors.Wrap(err, "failed to get unprocessed state hash")
		}

		_, err = p.stateStoreContract.GetStateInfoByState(nil, unprocessedStateHash)
		if err != nil {
			if isStateDoesntExist(err) {
				err = p.setStatusFailed("", err.Error(), &unprocessedState)
				if err != nil {
					return errors.Wrap(err, "failed to set status failed in db")
				}

				continue
			}
			return errors.Wrap(err, "failed to retrieve root hash from contract")
		}

		unprocessedState.Status = data.StatusCompleted
		err = p.state.DB.CommittedStatesQ().Update(&unprocessedState)
		if err != nil {
			return errors.Wrap(err, "failed to update unprocessed state status to complete")
		}
	}

	return nil
}

func getStateHash(state data.CommittedState) (*big.Int, error) {
	var claimsTreeHash merkletree.Hash
	var revocationsTreeHash merkletree.Hash
	var rootsTreeHash merkletree.Hash

	copy(claimsTreeHash[:], state.ClaimsTreeRoot)
	copy(revocationsTreeHash[:], state.RevocationsTreeRoot)
	copy(rootsTreeHash[:], state.RootsTreeRoot)

	hash, err := merkletree.HashElems(
		claimsTreeHash.BigInt(),
		revocationsTreeHash.BigInt(),
		rootsTreeHash.BigInt(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate state hash")
	}

	return hash.BigInt(), nil
}

func isStateDoesntExist(err error) bool {
	return strings.Contains(err.Error(), StateDoesntExistErrMsg)
}
