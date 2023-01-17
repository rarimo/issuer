package issuer

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	identityPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity"
)

const (
	ClaimIssueCallBackPath = "/integrations/qid-issuer/v1/claims/offers/callback"
)

var (
	ErrClaimIsNotExist               = errors.New("claim is not exist")
	ErrClaimRetrieverIsNotClaimOwner = errors.New("claim retriever is not claim owner")
	ErrMessageRecipientIsNotIssuer   = errors.New("the message recipient is not an issuer")
)

type issuer struct {
	*identityPkg.Identity
	schemaBuilder *schemas.Builder
	claimsOffersQ data.ClaimsOffersQ
	domain        string
}
