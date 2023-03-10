package publisher

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

func (p *publisher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case publishStateInfo := <-p.pendingQueue:
			running.UntilSuccess(ctx, p.log, statePublisherRunnerName,
				func(ctx context.Context) (bool, error) {
					return true, p.runner(ctx, publishStateInfo)
				}, p.runnerPeriod, p.runnerPeriod,
			)
		}
	}
}

func (p *publisher) runner(ctx context.Context, publishStateInfo *publishedStateInfo) error {
	receipt, err := p.waitMined(ctx, publishStateInfo.Tx)
	if err != nil {
		if errors.Is(err, ErrTransactionFailed) {
			p.log.WithField("reason", err.Error()).Info("Failed to wait mined tx")
			err = p.setStatusFailed(publishStateInfo.Tx.Hash().Hex(), err.Error(), publishStateInfo.CommittedState)
			if err != nil {
				return errors.Wrap(err, "failed to set status failed in db")
			}
			return nil
		}

		err = p.setStatusFailed(publishStateInfo.Tx.Hash().Hex(), err.Error(), publishStateInfo.CommittedState)
		if err != nil {
			return errors.Wrap(err, "failed to set status failed in db")
		}

		return errors.Wrap(err, "failed to wait mined tx")
	}

	err = p.setStatusCompleted(ctx, receipt, publishStateInfo.CommittedState)
	if err != nil {
		return errors.Wrap(err, "failed to set status completed")
	}

	p.log.WithField("tx_hash", publishStateInfo.Tx.Hash().Hex()).Info("State was successfully transited")
	return nil
}

func (p *publisher) setStatusFailed(txHash, reason string, committedState *data.CommittedState) error {
	committedState.TxID = txHash
	committedState.Status = data.StatusFailed
	committedState.Message = reason

	err := p.committedStatesQ.Update(committedState)
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

	err = p.committedStatesQ.Update(committedState)
	if err != nil {
		return errors.Wrap(err, "failed to update committed state in db")
	}

	return nil
}
