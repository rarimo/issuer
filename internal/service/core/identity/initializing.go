package identity

import (
	"context"
	"time"

	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (iden *Identity) generateNewIdentity(ctx context.Context) error {
	iden.log.Info("Generating the new Identity")

	if iden.babyJubJubPrivateKey == nil {
		return errors.New("error generating new identity, babyJubJubPrivateKey is nil")
	}

	identifier, authClaim, err := iden.State.SetupGenesis(iden.babyJubJubPrivateKey.Public())
	if err != nil {
		return errors.Wrap(err, "failed to setup genesis state")
	}

	iden.Identifier = identifier
	iden.log.
		WithField("identifier", identifier.String()).
		Infof("The new Identity successfully generated")

	if err := iden.saveAuthClaimModel(ctx, authClaim); err != nil {
		return errors.Wrap(err, "failed to save auth claim to db")
	}

	err = iden.State.CommittedStateQ.Insert((&state.CommittedState{
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

func (iden *Identity) saveAuthClaimModel(ctx context.Context, coreAuthClaim *core.Claim) error {
	if iden.babyJubJubPrivateKey == nil {
		return errors.New("error generating new identity, babyJubJubPrivateKey is nil")
	}

	authClaimData, err := claims.GenerateAuthClaimData(iden.babyJubJubPrivateKey.Public())
	if err != nil {
		return errors.Wrap(err, "failed to generate auth claim data")
	}

	authClaim := &data.Claim{
		CoreClaim:  data.NewCoreClaim(coreAuthClaim),
		SchemaURL:  schemas.AuthBJJCredentialSchemaURL,
		SchemaType: claims.AuthBJJCredentialClaimType,
		Data:       authClaimData,
	}

	authClaim.MTP, err = iden.GenerateMTP(ctx, coreAuthClaim)
	if err != nil {
		return errors.Wrap(err, "failed to generate proof")
	}

	err = iden.State.ClaimsQ.Insert(authClaim)
	if err != nil {
		return errors.Wrap(err, "failed to insert the auth claim to the db")
	}

	iden.AuthClaim = authClaim

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

	iden.Identifier, err = core.IdGenesisFromIdenState(core.TypeDefault, genesisStateHash.BigInt())
	if err != nil {
		return errors.Wrap(err, "failed to generate identifier from the genesis state")
	}

	if authClaim == nil || authClaim.CoreClaim == nil || authClaim.CoreClaim.Claim == nil {
		return errors.New("auth claim is nil")
	}

	authClaim.MTP, err = iden.GenerateMTP(ctx, authClaim.CoreClaim.Claim)
	if err != nil {
		return errors.Wrap(err, "failed to generate proof")
	}

	iden.AuthClaim = authClaim

	return nil
}
