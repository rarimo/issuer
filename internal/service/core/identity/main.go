package identity

import (
	"context"
	"math/big"

	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/pkg/errors"

	"gitlab.com/q-dev/q-id/issuer/internal/config"
	statePkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state_publisher"
)

func New(ctx context.Context, cfg config.Config) (*Identity, error) {
	state, err := statePkg.NewIdentityState(ctx, statePkg.Config{
		DB:             cfg.DB(),
		IdentityConfig: cfg.Identity(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize the identity state")
	}

	statePublisher, err := statepublisher.New(&statepublisher.Config{
		Log:            cfg.Log(),
		EthConfig:      cfg.EthClient(),
		StatePublisher: cfg.StatePublisher(),
	}, state)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize the state publisher")
	}

	identity := &Identity{
		log:                  cfg.Log().WithField("service", "Identity"),
		State:                state,
		babyJubJubPrivateKey: cfg.Identity().BabyJubJubPrivateKey,
		circuitsPath:         cfg.Identity().CircuitsPath,
	}

	err = identity.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init Identity")
	}

	identity.State.SetIdentityInfo(&statePkg.IdentityInfo{
		BabyJubJubPrivateKey: identity.babyJubJubPrivateKey,
		Identifier:           &identity.Identifier.ID,
		AuthClaim:            identity.AuthClaim.CoreClaim.Claim,
	})

	go statePublisher.Run(ctx)

	return identity, nil
}

func (iden *Identity) Init(ctx context.Context) error {
	genesisStateRaw, err := iden.State.DB.CommittedStatesQ().GetGenesis()
	if err != nil {
		return errors.Wrap(err, "failed to get genesis state")
	}

	authClaim, err := iden.State.DB.ClaimsQ().GetAuthClaim()
	if err != nil {
		return errors.Wrap(err, "failed to get auth claim")
	}

	if (genesisStateRaw == nil && authClaim != nil) || (genesisStateRaw != nil && authClaim == nil) {
		return errors.New("only one of the genesis state or auth claim is exist")
	}

	if genesisStateRaw == nil || authClaim == nil {
		iden.log.Info("Identity not found")
		err = iden.generateNewIdentity(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to generate new Identity")
		}

		return nil
	}

	iden.log.Info("Identity found")
	err = iden.parseIdentity(ctx, authClaim, genesisStateRaw)
	if err != nil {
		return errors.Wrap(err, "failed to parse Identity")
	}

	iden.State.SetIdentityInfo(&statePkg.IdentityInfo{
		BabyJubJubPrivateKey: iden.babyJubJubPrivateKey,
		Identifier:           &iden.Identifier.ID,
		AuthClaim:            iden.AuthClaim.CoreClaim.Claim,
	})

	return nil
}

func (iden *Identity) Sign(singMessage *big.Int) ([]byte, error) {
	if !utils.CheckBigIntInField(singMessage) {
		return nil, errors.New("data to signBytes is too large")
	}

	sig := iden.babyJubJubPrivateKey.SignPoseidon(singMessage).Compress()
	return sig[:], nil
}
