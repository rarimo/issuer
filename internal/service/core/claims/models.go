package claims

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/iden3/go-schema-processor/verifiable"
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
	MTPUpdateURL                 = "/integrations/issuer/v1/public/claims/mtp/update/"

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
	NaturalPersonSchemaType     ClaimSchemaType = "NaturalPerson"
	IdentityProvidersSchemaType ClaimSchemaType = "IdentityProviders"
)

const (
	NaturalPersonSchemaName     = "Natural Person"
	IdentityProvidersSchemaName = "Identity Providers"
)

const (
	NaturalPersonSchemaPath     = "/json/NaturalPerson.json"
	IdentityProvidersSchemaPath = "/json/IdentityProviders.json"
)

type ClaimDataParseFunc = func([]byte) ([]byte, error)

type ClaimData struct {
	ClaimSchemaURL        string
	ClaimSchemaName       string
	ClaimDataValidateFunc validation.RuleFunc
	ClaimDataParseFunc    ClaimDataParseFunc
}

var ClaimSchemaTypeList = map[string]ClaimSchemaType{
	NaturalPersonSchemaType.ToRaw():     NaturalPersonSchemaType,
	IdentityProvidersSchemaType.ToRaw(): IdentityProvidersSchemaType,
}

var ClaimSchemaList = map[ClaimSchemaType]ClaimData{
	NaturalPersonSchemaType: {
		ClaimSchemaURL:        NaturalPersonSchemaPath,
		ClaimSchemaName:       NaturalPersonSchemaName,
		ClaimDataValidateFunc: validationPkg.MustBeNaturalPersonCredentials,
		ClaimDataParseFunc:    validationPkg.ParseNaturalPersonCredentials,
	},
	IdentityProvidersSchemaType: {
		ClaimSchemaURL:        IdentityProvidersSchemaPath,
		ClaimSchemaName:       IdentityProvidersSchemaName,
		ClaimDataValidateFunc: validationPkg.MustBeIdentityProvidersCredentials,
		ClaimDataParseFunc:    validationPkg.ParseIdentityProvidersCredentials,
	},
}

type Iden3SparseMerkleTreeProofWithID struct {
	verifiable.Iden3SparseMerkleTreeProof

	ID string `json:"id"`
}
