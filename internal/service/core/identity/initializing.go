package identity

import (
	"context"
	"time"

	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

func (iden *Identity) generateNewIdentity(ctx context.Context, cfg *config.IdentityConfig) error {
	iden.log.Info("Generating the new Identity")

	identifier, authClaim, err := iden.State.SetupGenesis(cfg.BabyJubJubPrivateKey.Public())
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
	authClaimData, err := claims.GenerateAuthClaimData(iden.babyJubJubPrivateKey.Public())
	if err != nil {
		return errors.Wrap(err, "failed to generate auth claim data")
	}

	authClaim := &data.Claim{
		CoreClaim:  data.NewCoreClaim(coreAuthClaim),
		SchemaURL:  schemas.AuthBJJCredentialSchemaURL,
		SchemaType: claims.AuthBJJCredentialSchemaType.ToRaw(),
		Data:       authClaimData,
	}

	authClaim.MTP, err = iden.generateProof(ctx, coreAuthClaim, nil)
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
	genesisStateHash, err := state.CommittedStateFromRaw(genesisStateRaw).StateHash()
	if err != nil {
		return errors.Wrap(err, "failed to get genesis state hash")
	}

	iden.Identifier, err = core.IdGenesisFromIdenState(core.TypeDefault, genesisStateHash.BigInt())
	if err != nil {
		return errors.Wrap(err, "failed to generate identifier from the genesis state")
	}

	authClaim.MTP, err = iden.generateProof(ctx, authClaim.CoreClaim.Claim, nil)
	if err != nil {
		return errors.Wrap(err, "failed to generate proof")
	}

	iden.AuthClaim = authClaim

	return nil
}
