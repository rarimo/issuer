package issuer

import (
	"github.com/iden3/go-merkletree-sql/v2"
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

const (
	uuidStringSize = 36
)

var (
	ErrProofVerifyFailed             = errors.New("failed to verify proof")
	ErrClaimIsNotExist               = errors.New("claim is not exist")
	ErrClaimOfferIsNotExist          = errors.New("claim offer is not exist")
	ErrClaimRetrieverIsNotClaimOwner = errors.New("claim retriever is not claim owner")
	ErrMessageRecipientIsNotIssuer   = errors.New("the message recipient is not an issuer")
	ErrRepeatedCallbackRequest       = errors.New("repeated callback request")
	ErrClaimIsAlreadyRevoked         = errors.New("claim is already revoked")
	ErrInvalidCredentialID           = errors.New("invalid credential id")
)

type issuer struct {
	*identityPkg.Identity
	schemaBuilder       *schemas.Builder
	claimsOffersQ       data.ClaimsOffersQ
	authVerificationKey []byte
	baseURL             string
}

// ClaimInclusionMTP info required to check that claim is included in the issuer's claims tree.
// It is required for the extended Iden3 protocol tha supports cross chain with Rarimo.
type ClaimInclusionMTP struct {
	Issuer struct {
		State              *string `json:"state,omitempty"`
		RootOfRoots        *string `json:"rootOfRoots,omitempty"`
		ClaimsTreeRoot     *string `json:"claimsTreeRoot,omitempty"`
		RevocationTreeRoot *string `json:"revocationTreeRoot,omitempty"`
	} `json:"issuer"`
	MTP merkletree.Proof `json:"mtp"`
}
