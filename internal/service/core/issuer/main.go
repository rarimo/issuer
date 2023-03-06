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
	PublishStateOnChain(context.Context) (string, error)
	GetIdentifier() string
	CreateClaimOffer(*core.DID, string) (*protocol.CredentialsOfferMessage, error)
	IssueClaim(context.Context, *core.DID, *time.Time, resources.ClaimSchemaType, []byte) (string, error)
	OfferCallback(context.Context, *requests.OfferCallbackRequest) (*protocol.CredentialIssuanceMessage, error)
	GetRevocationStatus(context.Context, *big.Int) (*verifiable.RevocationStatus, error)
	RevokeClaim(context.Context, *core.ID, resources.ClaimSchemaType) error
}

func New(ctx context.Context, cfg config.Config) (Issuer, error) {
	identity, err := identityPkg.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the new identity")
	}

	schemaBuilder, err := schemas.NewBuilder(ctx)
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
		baseURL:             cfg.Issuer().BaseURL,
	}, nil
}
