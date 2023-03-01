package publisher

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
	"gitlab.com/q-dev/q-id/issuer/internal/data/pg"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/publisher/contracts"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/zkp"
)

type Publisher interface {
	PublishState(ctx context.Context, stInfo *StateTransitionInfo, committedState *data.CommittedState) (string, error)
	Run(ctx context.Context)
}

func NewPublisher(cfg *Config) (Publisher, error) {
	stateStoreContract, err := contracts.NewStateStore(*cfg.EthConfig.StateStorageContract, cfg.EthConfig.EthClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new state store contract")
	}

	chainID, err := cfg.EthConfig.EthClient.ChainID(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chainID")
	}

	publisherInstance := &publisher{
		log:                cfg.Log,
		committedStatesQ:   pg.NewCommittedStateQ(cfg.DB),
		ethClient:          cfg.EthConfig.EthClient,
		stateStoreContract: stateStoreContract,
		privateKey:         cfg.EthConfig.PrivateKey,
		address:            crypto.PubkeyToAddress(cfg.EthConfig.PrivateKey.PublicKey),
		chainID:            chainID,
		pendingQueue:       make(chan *publishedStateInfo),
	}

	err = publisherInstance.processPreviousSessionCommit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact pending queue")
	}

	return publisherInstance, nil
}

func (p *publisher) PublishState(
	ctx context.Context,
	tsInfo *StateTransitionInfo,
	committedState *data.CommittedState,
) (string, error) {
	zkpArgs, err := parseZKPArgs(tsInfo.ZKProof)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse SKProof to contract readable format")
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
			return errors.Wrap(err, "failed to call transit state contract method")
		}

		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to call retryChainCall")
	}

	// Runner will wait while transaction will be mined
	select {
	case p.pendingQueue <- &publishedStateInfo{CommittedState: committedState, Tx: tx}:
	default:
		return "", nil
	}

	return tx.Hash().Hex(), nil
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
	unprocessedStateList, err := p.committedStatesQ.New().WhereStatus(data.StatusProcessing).Select()
	if err != nil {
		return errors.Wrap(err, "failed to select unprocessed committed states from db")
	}

	for i, unprocessedState := range unprocessedStateList {
		unprocessedStateHash, err := getStateHash(unprocessedState)
		if err != nil {
			return errors.Wrap(err, "failed to get unprocessed state hash")
		}

		_, err = p.stateStoreContract.GetStateInfoByState(nil, unprocessedStateHash)
		if err != nil {
			if isStateDoesntExist(err) {
				unprocessedState.Status = data.StatusFailed
				err = p.committedStatesQ.Update(&unprocessedStateList[i])
				if err != nil {
					return errors.Wrap(err, "failed to update unprocessed state status to failed")
				}

				continue
			}
			return errors.Wrap(err, "failed to retrieve root hash from contract")
		}

		unprocessedState.Status = data.StatusCompleted
		err = p.committedStatesQ.Update(&unprocessedStateList[i])
		if err != nil {
			return errors.Wrap(err, "failed to update unprocessed state status to complete")
		}
	}

	return nil
}

func getStateHash(state data.CommittedState) (*big.Int, error) {
	var claimsTreeHash *merkletree.Hash
	var revocationsTreeHash *merkletree.Hash
	var rootsTreeHash *merkletree.Hash

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
