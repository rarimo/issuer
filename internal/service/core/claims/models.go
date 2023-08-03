package claims

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validationPkg "github.com/rarimo/issuer/internal/service/core/claims/validation"
)

const (
	CredentialStatusCheckURL = "/integrations/issuer/v1/public/claims/revocations/check/"
	MTPUpdateURL             = "/integrations/issuer/v1/public/claims/mtp-update/"

	AuthBJJCredentialClaimType = "AuthBJJCredential" //nolint
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
