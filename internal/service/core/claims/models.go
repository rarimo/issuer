package claims

import (
	"time"

	"github.com/iden3/go-schema-processor/processor"
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
	credentialStatusCheckURL     = "/integrations/qid-issuer/v1/claims/revocations/check/"

	Iden3CredentialSchemaType   = "Iden3Credential"   //nolint
	AuthBJJCredentialSchemaType = "AuthBJJCredential" //nolint

	AuthBJJCredentialHash      = "ca938857241db9451ea329256b9c06e5"                                                             //nolint
	AuthBJJCredentialSchemaURL = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld" //nolint
)

var (
	ErrIncorrectSchemaHashLength = errors.New("incorrect schema hash bytes length")
	ErrIncorrectSubjectPosition  = errors.New("incorrect subject position")
	ErrIDPositionIsNotSpecified  = errors.New("id position is not specified")
)

type CoreClaimData struct {
	SchemaHash      string
	Slots           processor.ParsedSlots
	SubjectID       string
	Expiration      time.Time
	SubjectPosition string
}
