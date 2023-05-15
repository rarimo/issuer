package identity

import (
	"context"
	"time"

	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (iden *Identity) generateNewIdentity(ctx context.Context) error {
	iden.log.Info("Generating the new Identity")

	if iden.babyJubJubPrivateKey == nil {
		return errors.New("error generating new identity, babyJubJubPrivateKey is nil")
	}

	did, authClaim, err := iden.State.SetupGenesis(iden.babyJubJubPrivateKey.Public())
	if err != nil {
		return errors.Wrap(err, "failed to setup genesis state")
	}

	iden.Identifier = did
	iden.log.
		WithField("did", did.String()).
		Infof("The new Identity successfully generated")

	if err := iden.saveAuthClaimModel(authClaim); err != nil {
		return errors.Wrap(err, "failed to save auth claim to db")
	}

	err = iden.State.DB.CommittedStatesQ().Insert((&state.CommittedState{
		Status:              data.StatusCompleted,
		CommitInfo:          nil,
		CreatedAt:           time.Now(),
		IsGenesis:           true,
		RootsTreeRoot:       iden.State.RootsTree.Root(),
		ClaimsTreeRoot:      iden.State.ClaimsTree.Root(),
		RevocationsTreeRoot: iden.State.RevocationsTree.Root(),
	}).ToRaw())
	if err != nil {
		return errors.Wrap(err, "failed to insert genesis state to committed states")
	}

	return nil
}

func (iden *Identity) parseIdentity(
	ctx context.Context,
	authClaim *data.Claim,
	genesisStateRaw *data.CommittedState,
) error {
	genesisState, err := state.CommittedStateFromRaw(genesisStateRaw)
	if err != nil {
		return errors.Wrap(err, "failed to parse genesis state")
	}

	genesisStateHash, err := genesisState.StateHash()
	if err != nil {
		return errors.Wrap(err, "failed to get genesis state hash")
	}

	didType, err := core.BuildDIDType(core.DIDMethodIden3, core.NoChain, core.NoNetwork)
	if err != nil {
		return errors.Wrap(err, "failed to build did type")
	}

	iden.Identifier, err = core.DIDGenesisFromIdenState(didType, genesisStateHash.BigInt())
	if err != nil {
		return errors.Wrap(err, "failed to generate identifier from the genesis state")
	}

	if authClaim == nil || authClaim.CoreClaim == nil || authClaim.CoreClaim.Claim == nil {
		return errors.New("auth claim is nil")
	}

	iden.AuthClaim = authClaim

	iden.log.
		WithField("did", iden.Identifier.String()).
		Infof("The Identity successfully loaded")

	return nil
}

func (iden *Identity) saveAuthClaimModel(coreAuthClaim *core.Claim) error {
	if iden.babyJubJubPrivateKey == nil {
		return errors.New("error generating new identity, babyJubJubPrivateKey is nil")
	}

	authClaimData, err := claims.GenerateAuthClaimData(iden.babyJubJubPrivateKey.Public())
	if err != nil {
		return errors.Wrap(err, "failed to generate auth claim data")
	}

	authClaim := &data.Claim{
		CoreClaim:  data.NewCoreClaim(coreAuthClaim),
		ClaimType:  claims.AuthBJJCredentialClaimType,
		Credential: authClaimData,
	}

	err = iden.State.DB.ClaimsQ().Insert(authClaim)
	if err != nil {
		return errors.Wrap(err, "failed to insert the auth claim to the db")
	}

	iden.AuthClaim = authClaim

	return nil
}
