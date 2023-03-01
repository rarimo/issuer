package claims

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By the default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"

	BabyJubSignatureType = "BJJSignature2021"

	CorrectSchemaHashBytesLength = 16
	ExpirationWithoutExpiration  = 0
	CredentialStatusCheckURL     = "/integrations/qid-issuer/v1/claims/revocations/check/"

	AuthBJJCredentialClaimType = "AuthBJJCredential" //nolint
)

const (
	IssuanceDateCredentialField = "issuanceDate"
)

var (
	ErrIDPositionIsNotSpecified = errors.New("id position is not specified")
)
