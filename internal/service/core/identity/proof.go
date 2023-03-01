package identity

import (
	"context"
	"encoding/json"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (iden *Identity) GenerateMTP(
	ctx context.Context,
	claim *core.Claim,
) ([]byte, error) {
	if claim == nil {
		return nil, errors.New("failed to generate proof, claim is nil")
	}

	lastCommittedState, lastCommittedStateRaw, err := iden.GetLatestState()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state")
	}

	inclusionProof, err := iden.State.GetInclusionProof(ctx, claim, lastCommittedState.ClaimsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate inclusion proof")
	}
	if !inclusionProof.Existence {
		return nil, ErrClaimWasNotPublishedYet
	}

	lastCommittedStateHash, err := lastCommittedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state hash")
	}

	mtProof := verifiable.Iden3SparseMerkleProof{
		Type: verifiable.Iden3SparseMerkleProofType,
		MTP:  inclusionProof,
		IssuerData: verifiable.IssuerData{
			ID: iden.Identifier.String(),
			State: verifiable.State{
				RootOfRoots:        strToPtr(lastCommittedState.RootsTreeRoot.Hex()),
				ClaimsTreeRoot:     strToPtr(lastCommittedState.ClaimsTreeRoot.Hex()),
				RevocationTreeRoot: strToPtr(lastCommittedState.RevocationsTreeRoot.Hex()),
				Value:              strToPtr(lastCommittedStateHash.Hex()),
			},
			MTP: inclusionProof,
		},
	}

	if lastCommittedStateRaw != nil && lastCommittedStateRaw.TxID != "" {
		blockTimestamp := int(lastCommittedStateRaw.BlockTimestamp)
		blockNumber := int(lastCommittedStateRaw.BlockNumber)
		mtProof.IssuerData.State.TxID = &lastCommittedStateRaw.TxID
		mtProof.IssuerData.State.BlockTimestamp = &blockTimestamp
		mtProof.IssuerData.State.BlockNumber = &blockNumber
	}

	proofRaw, err := json.Marshal(mtProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merkle tree proof")
	}

	return proofRaw, nil
}

func (iden *Identity) GetLatestState() (*state.CommittedState, *data.CommittedState, error) {
	lastCommittedStateRaw, err := iden.State.CommittedStateQ.WhereStatus(data.StatusCompleted).GetLatest()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get last committed state from db")
	}

	if lastCommittedStateRaw != nil {
		lastCommittedState, err := state.CommittedStateFromRaw(lastCommittedStateRaw)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to get latest committed state")
		}

		return lastCommittedState, lastCommittedStateRaw, nil
	}

	return &state.CommittedState{
		ClaimsTreeRoot:      iden.State.ClaimsTree.Root(),
		RevocationsTreeRoot: iden.State.RevocationsTree.Root(),
		RootsTreeRoot:       iden.State.RootsTree.Root(),
	}, nil, nil
}

func strToPtr(s string) *string {
	return &s
}
