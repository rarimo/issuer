package claims

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"

	validationPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/claims/validation"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By the default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"

	BabyJubSignatureType = "BJJSignature2021"

	CorrectSchemaHashBytesLength = 16
	ExpirationWithoutExpiration  = 0
	CredentialStatusCheckURL     = "/integrations/issuer/v1/claims/revocations/check/"

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
	DAOMembershipSchemaType ClaimSchemaType = "DAOMembership" //nolint
)

const (
	DAOMembershipSchemaName = "DAO Membership" //nolint
)

const (
	DAOMembershipSchemaPath = "/json/DAOMembership.json" //nolint
)

type ClaimDataParseFunc = func([]byte) ([]byte, error)

type ClaimData struct {
	ClaimSchemaURL        string
	ClaimSchemaName       string
	ClaimDataValidateFunc validation.RuleFunc
	ClaimDataParseFunc    ClaimDataParseFunc
}

var ClaimSchemaTypeList = map[string]ClaimSchemaType{
	DAOMembershipSchemaType.ToRaw(): DAOMembershipSchemaType,
}

var ClaimSchemaList = map[ClaimSchemaType]ClaimData{
	DAOMembershipSchemaType: {
		ClaimSchemaURL:        DAOMembershipSchemaPath,
		ClaimSchemaName:       DAOMembershipSchemaName,
		ClaimDataValidateFunc: validationPkg.MustBeDAOMembership,
		ClaimDataParseFunc:    validationPkg.ParseDAOMembership,
	},
}
