package claims

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"

	validationPkg "gitlab.com/rarimo/identity/issuer/internal/service/core/claims/validation"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By the default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"

	BabyJubSignatureType = "BJJSignature2021"

	CorrectSchemaHashBytesLength = 16
	ExpirationWithoutExpiration  = 0
	CredentialStatusCheckURL     = "/integrations/issuer/v1/public/claims/revocations/check/"

	AuthBJJCredentialClaimType = "AuthBJJCredential" //nolint
)

var (
	ErrIDPositionIsNotSpecified = errors.New("id position is not specified")
)

type ClaimSchemaType string

func (c ClaimSchemaType) ToRaw() string {
	return string(c)
}

const (
	NaturalPersonSchemaType ClaimSchemaType = "NaturalPerson" //nolint
)

const (
	NaturalPersonSchemaName = "Natural Person" //nolint
)

const (
	NaturalPersonSchemaPath = "/json/NaturalPerson.json" //nolint
)

type ClaimDataParseFunc = func([]byte) ([]byte, error)

type ClaimData struct {
	ClaimSchemaURL        string
	ClaimSchemaName       string
	ClaimDataValidateFunc validation.RuleFunc
	ClaimDataParseFunc    ClaimDataParseFunc
}

var ClaimSchemaTypeList = map[string]ClaimSchemaType{
	NaturalPersonSchemaType.ToRaw(): NaturalPersonSchemaType,
}

var ClaimSchemaList = map[ClaimSchemaType]ClaimData{
	NaturalPersonSchemaType: {
		ClaimSchemaURL:        NaturalPersonSchemaPath,
		ClaimSchemaName:       NaturalPersonSchemaName,
		ClaimDataValidateFunc: validationPkg.MustBeNaturalPersonCredentials,
		ClaimDataParseFunc:    validationPkg.ParseNaturalPersonCredentials,
	},
}
