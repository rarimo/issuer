package claims

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iden3/go-schema-processor/processor"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/schemas"
	validatonPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/validation"
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

type ClaimSchemaType string

func (c ClaimSchemaType) ToRaw() string {
	return string(c)
}

// TODO: everything below should be moved in the some config or something else

const (
	Iden3CredentialSchema            ClaimSchemaType = "Iden3Credential"           //nolint
	AuthBJJCredentialSchemaType      ClaimSchemaType = "AuthBJJCredential"         //nolint
	KYCFullNameCredentialsSchemaType ClaimSchemaType = "KYC Full name credentials" //nolint
	QDAOMembershipSchemaType         ClaimSchemaType = "Q DAO Membership"
)

type ClaimDataParseFunc = func([]byte) ([]byte, error)

type KYCSchema struct {
	ClaimSchemaURL        string
	ClaimDataValidateFunc validation.RuleFunc
	ClaimDataParseFunc    ClaimDataParseFunc
}

var ClaimSchemaTypeList = map[string]ClaimSchemaType{
	"KYC Full name credentials": KYCFullNameCredentialsSchemaType,
	"Q DAO Membership":          QDAOMembershipSchemaType,
}

var ClaimSchemaList = map[ClaimSchemaType]KYCSchema{
	KYCFullNameCredentialsSchemaType: {
		ClaimSchemaURL:        schemas.KYCFullNameCredentialSchemaURL,
		ClaimDataValidateFunc: validatonPkg.MustBeKYCFullNameCredentials,
		ClaimDataParseFunc:    validatonPkg.ConvertKYCFullNameCredentials,
	},
	QDAOMembershipSchemaType: {
		ClaimSchemaURL:        schemas.QDAOMembershipSchemaURL,
		ClaimDataValidateFunc: validatonPkg.MustBeQDAOMembership,
		ClaimDataParseFunc:    validatonPkg.ConvertQDAOMembership,
	},
}
