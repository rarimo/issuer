package issuer

import (
	"context"
	"math/big"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/rarimo/identity/issuer/internal/config"
	"gitlab.com/rarimo/identity/issuer/internal/data/pg"
	"gitlab.com/rarimo/identity/issuer/internal/service/api/requests"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/claims"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/rarimo/identity/issuer/internal/service/core/identity"
	statePkg "gitlab.com/rarimo/identity/issuer/internal/service/core/identity/state"
)

type Issuer interface {
	GetIdentifier() string
	CreateClaimOffer(*core.DID, string) (*protocol.CredentialsOfferMessage, error)
	IssueClaim(context.Context, *core.DID, *time.Time, claims.ClaimSchemaType, []byte) (string, error)
	OfferCallback(context.Context, *requests.OfferCallbackRequest) (*protocol.CredentialIssuanceMessage, error)
	GetRevocationStatus(context.Context, *big.Int) (*verifiable.RevocationStatus, error)
	GetInclusionMTP(ctx context.Context, claimID uuid.UUID) (*ClaimInclusionMTP, error)
	RevokeClaim(context.Context, *core.ID, claims.ClaimSchemaType) error
}

func New(ctx context.Context, cfg config.Config) (Issuer, error) {
	identity, err := identityPkg.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the new identity")
	}

	schemaBuilder, err := schemas.NewBuilder(ctx, cfg.Issuer().SchemasBaseURL)
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
