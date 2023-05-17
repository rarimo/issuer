package statepublisher

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/running"

	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (p *publisher) Run(ctx context.Context) {
	fmt.Println(p.publishPeriod)
	ticker := time.NewTicker(p.publishPeriod)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			running.UntilSuccess(ctx, p.log, statePublisherRunnerName,
				func(ctx context.Context) (bool, error) {
					return true, p.runner(ctx)
				}, p.retryPeriod, p.retryPeriod,
			)

			ticker.Reset(p.publishPeriod)
		}
	}
}

func (p *publisher) runner(ctx context.Context) error {
	stateTransitionInfo, stateCommit, err := p.state.GenerateStateCommitment(ctx)
	if err != nil {
		if errors.Is(err, state.ErrStateWasntChanged) {
			return nil
		}

		return errors.Wrap(err, "failed to generate state commitment")
	}

	tx, err := p.sendTransaction(ctx, stateTransitionInfo, stateCommit)
	if err != nil {
		return errors.Wrap(err, "failed to send transaction")
	}

	err = p.waitTransaction(ctx, stateCommit, tx)
	if err != nil {
		return errors.Wrap(err, "failed to wait transaction")
	}

	return nil
}

func (p *publisher) setStatusFailed(txHash, reason string, committedState *data.CommittedState) error {
	committedState.TxID = txHash
	committedState.Status = data.StatusFailed
	committedState.Message = reason

	err := p.state.DB.CommittedStatesQ().Update(committedState)
	if err != nil {
		return errors.Wrap(err, "failed to update committed state in db")
	}

	return nil
}

func (p *publisher) setStatusCompleted(
	ctx context.Context,
	receipt *types.Receipt,
	committedState *data.CommittedState,
) error {
	block, err := p.ethClient.BlockByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return errors.Wrap(err, "failed to get ethereum block by number")
	}

	committedState.TxID = receipt.TxHash.Hex()
	committedState.Status = data.StatusCompleted
	committedState.BlockTimestamp = block.Time()
	committedState.BlockNumber = block.NumberU64()

	err = p.state.DB.CommittedStatesQ().Update(committedState)
	if err != nil {
		return errors.Wrap(err, "failed to update committed state in db")
	}

	return nil
}
