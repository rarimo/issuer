package identity

import (
	"context"
	"encoding/json"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (iden *Identity) GenerateProof(
	ctx context.Context,
	claim *core.Claim,
	claimTreeRoot *merkletree.Hash,
) ([]byte, error) {
	if claim == nil {
		return nil, errors.New("failed to generate proof, claim is nil")
	}

	inclusionProof, err := iden.State.GetInclusionProof(ctx, claim, claimTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate inclusion proof")
	}

	mtProof := verifiable.Iden3SparseMerkleProof{}
	mtProof.Type = verifiable.Iden3SparseMerkleProofType
	mtProof.MTP = inclusionProof

	stateHash, err := iden.State.GetCurrentStateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current state hash")
	}

	stateHashHex := stateHash.Hex()
	claimsRootHex := iden.State.ClaimsTree.Root().Hex()
	mtProof.IssuerData = verifiable.IssuerData{
		ID: iden.Identifier,
		State: verifiable.State{
			Value:          &stateHashHex,
			ClaimsTreeRoot: &claimsRootHex,
		},
		MTP: inclusionProof,
	}

	proofRaw, err := json.Marshal(mtProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merkle tree proof")
	}

	return proofRaw, nil
}

func (iden *Identity) GenerateMTP(
	ctx context.Context,
	claim *core.Claim,
) ([]byte, error) {
	lastCommittedStateRaw, err := iden.State.CommittedStateQ.WhereStatus(data.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last committed state from db")
	}

	if lastCommittedStateRaw.TxID == "" {
		return nil, errors.New("failed to generate full merkle tree proof txID is empty")
	}

	var claimTreeRoot merkletree.Hash
	copy(claimTreeRoot[:], lastCommittedStateRaw.ClaimsTreeRoot)
	inclusionProof, err := iden.State.GetInclusionProof(ctx, claim, &claimTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate merkle tree proof")
	}
	if !inclusionProof.Existence {
		return nil, ErrClaimWasNotPublishedYet
	}

	lastCommittedState, err := state.CommittedStateFromRaw(lastCommittedStateRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state")
	}

	lastCommittedStateHash, err := lastCommittedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state hash")
	}

	blockTimestamp := int(lastCommittedStateRaw.BlockTimestamp)
	blockNumber := int(lastCommittedStateRaw.BlockNumber)
	mtp := &verifiable.Iden3SparseMerkleProof{
		Type: verifiable.Iden3SparseMerkleProofType,
		IssuerData: verifiable.IssuerData{
			ID: iden.Identifier,
			State: verifiable.State{
				TxID:               &lastCommittedStateRaw.TxID,
				BlockTimestamp:     &blockTimestamp,
				BlockNumber:        &blockNumber,
				RootOfRoots:        strToPtr(lastCommittedState.RootsTreeRoot.Hex()),
				ClaimsTreeRoot:     strToPtr(lastCommittedState.ClaimsTreeRoot.Hex()),
				RevocationTreeRoot: strToPtr(lastCommittedState.RevocationsTreeRoot.Hex()),
				Value:              strToPtr(lastCommittedStateHash.Hex()),
			},
		},
		MTP: inclusionProof,
	}

	mtpRaw, err := json.Marshal(mtp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal merkle tree proof")
	}

	return mtpRaw, nil
}

func strToPtr(s string) *string {
	return &s
}
