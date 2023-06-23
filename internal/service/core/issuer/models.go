package issuer

import (
	"github.com/pkg/errors"

	"gitlab.com/rarimo/identity/issuer/internal/data"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/rarimo/identity/issuer/internal/service/core/identity"
)

const (
	ClaimIssueCallBackPath = "/integrations/issuer/v1/public/claims/offers/callback"
	GetClaimPath           = "/integrations/issuer/v1/private/claims/"
	basicAuthKeyPath       = "/auth/verification_key.json"
)

var (
	ErrProofVerifyFailed             = errors.New("failed to verify proof")
	ErrClaimIsNotExist               = errors.New("claim is not exist")
	ErrClaimOfferIsNotExist          = errors.New("claim offer is not exist")
	ErrClaimRetrieverIsNotClaimOwner = errors.New("claim retriever is not claim owner")
	ErrMessageRecipientIsNotIssuer   = errors.New("the message recipient is not an issuer")
	ErrRepeatedCallbackRequest       = errors.New("repeated callback request")
	ErrClaimIsAlreadyRevoked         = errors.New("claim is already revoked")
)

type issuer struct {
	*identityPkg.Identity
	schemaBuilder       *schemas.Builder
	claimsOffersQ       data.ClaimsOffersQ
	authVerificationKey []byte
	baseURL             string
}
