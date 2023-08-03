package identity

import (
	"context"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/rarimo/issuer/internal/data"
	"github.com/rarimo/issuer/internal/service/core/claims"
	"github.com/rarimo/issuer/internal/service/core/identity/state"
)

// GenerateIncMTProof generates merkletree proof for inclusion claim in issuer's claims tree.
func (iden *Identity) GenerateIncMTProof(
	ctx context.Context,
	claim *core.Claim,
	updateURL string, // contains URL for MTProof updating
	issuerData verifiable.IssuerData,
) (*claims.Iden3SparseMerkleTreeProof, error) {
	if claim == nil {
		return nil, errors.New("failed to generate proof, claim is nil")
	}

	lastCommittedState, _, err := iden.GetLatestState()
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

	coreClaimHex, err := claim.Hex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse core claim hex")
	}

	return &claims.Iden3SparseMerkleTreeProof{
		Type:       verifiable.Iden3SparseMerkleTreeProofType,
		MTP:        inclusionProof,
		IssuerData: issuerData,
		CoreClaim:  coreClaimHex,
		ID:         updateURL,
	}, nil
}

// GenerateIncSigProof generates signature proof for inclusion claim in issuer's claims tree.
func (iden *Identity) GenerateIncSigProof(
	claim *core.Claim,
	issuerData verifiable.IssuerData,
	issuerProofUpdateURL string,
) (*claims.BJJSignatureProof2021, error) {
	signature, err := claims.SignClaimEntry(claim, iden.Sign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign core claim")
	}

	coreClaimHex, err := claim.Hex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hex from auth core claim")
	}

	return &claims.BJJSignatureProof2021{
		Type:                 verifiable.BJJSignatureProofType,
		Signature:            signature,
		CoreClaim:            coreClaimHex,
		IssuerData:           issuerData,
		IssuerProofUpdateURL: issuerProofUpdateURL,
	}, nil
}

func (iden *Identity) CompactIssuerData(ctx context.Context, checkRevLink string) (*verifiable.IssuerData, error) {
	lastCommittedState, lastCommittedStateRaw, err := iden.GetLatestState()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state")
	}

	lastCommittedStateHash, err := lastCommittedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state hash")
	}

	authCoreClaimHex, err := iden.AuthClaim.CoreClaim.Claim.Hex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hex from auth core claim")
	}

	authInclusionProof, err := iden.State.GetInclusionProof(
		ctx,
		iden.AuthClaim.CoreClaim.Claim,
		lastCommittedState.ClaimsTreeRoot,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate inclusion proof")
	}
	if !authInclusionProof.Existence {
		return nil, ErrClaimWasNotPublishedYet
	}

	issuerData := verifiable.IssuerData{
		ID: iden.Identifier.String(),
		State: verifiable.State{
			RootOfRoots:        strToPtr(lastCommittedState.RootsTreeRoot.Hex()),
			ClaimsTreeRoot:     strToPtr(lastCommittedState.ClaimsTreeRoot.Hex()),
			RevocationTreeRoot: strToPtr(lastCommittedState.RevocationsTreeRoot.Hex()),
			Value:              strToPtr(lastCommittedStateHash.Hex()),
		},
		AuthCoreClaim: authCoreClaimHex,
		MTP:           authInclusionProof,
		CredentialStatus: &verifiable.CredentialStatus{
			ID:              checkRevLink,
			Type:            verifiable.SparseMerkleTreeProof,
			RevocationNonce: iden.AuthClaim.CoreClaim.GetRevocationNonce(),
		},
	}

	if lastCommittedStateRaw != nil && lastCommittedStateRaw.TxID != "" {
		blockTimestamp := int(lastCommittedStateRaw.BlockTimestamp)
		blockNumber := int(lastCommittedStateRaw.BlockNumber)
		issuerData.State.TxID = &lastCommittedStateRaw.TxID
		issuerData.State.BlockTimestamp = &blockTimestamp
		issuerData.State.BlockNumber = &blockNumber
	}

	return &issuerData, nil
}

func (iden *Identity) GetLatestState() (*state.CommittedState, *data.CommittedState, error) {
	lastCommittedStateRaw, err := iden.State.DB.CommittedStatesQ().WhereStatus(data.StatusCompleted).GetLatest()
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
