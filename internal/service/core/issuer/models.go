package issuer

import (
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity"
)

const (
	ClaimIssueCallBackPath = "/integrations/qid-issuer/v1/claims/offers/callback"
	basicAuthKeyPath       = "/key/auth.json"
)

var (
	ErrProofVerifyFailed             = errors.New("failed to verify proof")
	ErrClaimIsNotExist               = errors.New("claim is not exist")
	ErrClaimOfferIsNotExist          = errors.New("claim offer is not exist")
	ErrClaimRetrieverIsNotClaimOwner = errors.New("claim retriever is not claim owner")
	ErrMessageRecipientIsNotIssuer   = errors.New("the message recipient is not an issuer")
	ErrRepeatedCallbackRequest       = errors.New("repeated callback request")
)

type issuer struct {
	*identityPkg.Identity
	schemaBuilder       *schemas.Builder
	claimsOffersQ       data.ClaimsOffersQ
	authVerificationKey []byte
	baseURL             string
}
