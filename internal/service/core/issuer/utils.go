package issuer

import (
	"context"
	"strconv"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

func (isr *issuer) generateMTPForOfferCallback(ctx context.Context, claim *data.Claim) error {
	lastState, err := isr.State.CommittedStateQ.GetLatest()
	if err != nil {
		return errors.Wrap(err, "failed to get last committed state from db")
	}

	if !lastState.IsGenesis {
		claim.MTP, err = isr.GenerateMTP(ctx, claim.CoreClaim.Claim)
		if err != nil {
			return errors.Wrap(err, "failed to generate merkle tree proof")
		}
	}

	return nil
}

func (isr *issuer) retrieveClaim(idRaw string) (*data.Claim, error) {
	claimID, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse claim id")
	}

	claim, err := isr.State.ClaimsQ.Get(claimID)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return claim, nil
}

func checkClaimRetriever(claim *data.Claim, claimRetriever string) (bool, error) {
	claimRecipientID, err := claim.CoreClaim.GetID()
	if err != nil {
		return false, errors.Wrap(err, "failed to get claim recipient identifier")
	}

	if claimRecipientID.String() != claimRetriever {
		return false, nil
	}

	return true, nil
}
