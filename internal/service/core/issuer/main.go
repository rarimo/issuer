package issuer

import (
	"context"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/data/pg"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity"
)

type Issuer interface {
	PublishStateOnChain(ctx context.Context) (string, error)
	CreateClaimOffer(userID *core.ID, claimID uint64) (*protocol.CredentialsOfferMessage, error)
	IssueClaim(
		ctx context.Context, userID *core.ID, expiration time.Time, schemaType claims.ClaimSchemaType, schemaData []byte,
	) (uint64, error)
	OfferCallback(
		ctx context.Context, fetchMessage *protocol.CredentialFetchRequestMessage,
	) (*protocol.CredentialIssuanceMessage, error)
}

func New(ctx context.Context, cfg config.Config) (Issuer, error) {
	identity, err := identityPkg.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the new identity")
	}

	schemaBuilder, err := schemas.NewBuilder(cfg.Issuer().IpfsNodeURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new schema builder")
	}

	return &issuer{
		Identity:      identity,
		schemaBuilder: schemaBuilder,
		claimsOffersQ: pg.NewClaimsOffersQ(cfg.DB()),
		domain:        cfg.Issuer().Domain,
	}, nil
}
