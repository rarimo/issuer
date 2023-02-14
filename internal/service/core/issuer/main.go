package issuer

import (
	"context"
	"math/big"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/data/pg"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api/requests"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity"
	statePkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	resources "gitlab.com/q-dev/q-id/resources/claim_resources"
)

type Issuer interface {
	PublishStateOnChain(ctx context.Context) (string, error)
	GetIdentifier() string
	CreateClaimOffer(userID *core.ID, claimID string) (*protocol.CredentialsOfferMessage, error)
	IssueClaim(
		ctx context.Context, userID *core.ID, expiration time.Time,
		schemaType resources.ClaimSchemaType, schemaData []byte,
	) (uint64, error)
	OfferCallback(
		ctx context.Context, request *requests.OfferCallbackRequest,
	) (*protocol.CredentialIssuanceMessage, error)
	GetRevocationStatus(
		ctx context.Context, revID *big.Int,
	) (*verifiable.RevocationStatus, error)
	RevokeClaim(
		ctx context.Context, userID *core.ID, schemaType resources.ClaimSchemaType,
	) error
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

	authVerificationKey, err := statePkg.ReadFileByPath(cfg.Identity().CircuitsPath, basicAuthKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read basic auth verification key")
	}

	return &issuer{
		Identity:            identity,
		schemaBuilder:       schemaBuilder,
		claimsOffersQ:       pg.NewClaimsOffersQ(cfg.DB()),
		authVerificationKey: authVerificationKey,
		domain:              cfg.Issuer().Domain,
	}, nil
}
