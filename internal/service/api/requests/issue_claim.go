package requests

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/rarimo/identity/issuer/internal/service/core/claims"
	"gitlab.com/rarimo/identity/issuer/resources"
)

type IssueClaimRequest struct {
	UserDID    *core.DID
	ClaimType  claims.ClaimSchemaType
	Expiration *time.Time
	Credential []byte
}

type issueClaimRequestRaw struct {
	UserID    string
	ClaimType string
	Body      resources.IssueClaimRequest
}

func NewIssueClaim(r *http.Request, issuerID string) (*IssueClaimRequest, error) {
	requestBody := resources.IssueClaimRequest{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, errors.Wrap(err, "failed to decode json request body")
	}

	requestRaw := issueClaimRequestRaw{
		UserID:    chi.URLParam(r, UserIDPathParam),
		ClaimType: chi.URLParam(r, claimTypePathParam),
		Body:      requestBody,
	}

	if err := requestRaw.validate(issuerID); err != nil {
		return nil, err
	}

	schemaType := claims.ClaimSchemaTypeList[requestRaw.ClaimType]
	if err := validation.Validate(
		requestRaw.Body.Data.Attributes.CredentialSubject,
		validation.By(
			claims.ClaimSchemaList[schemaType].ClaimDataValidateFunc,
		),
	); err != nil {
		return nil, errors.Wrap(err, "invalid schema data")
	}

	parseData, err := claims.ClaimSchemaList[schemaType].ClaimDataParseFunc(
		requestRaw.Body.Data.Attributes.CredentialSubject,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse kyc full name data")
	}

	requestRaw.Body.Data.Attributes.CredentialSubject = parseData

	return requestRaw.parse(), nil
}

// nolint
func (req *issueClaimRequestRaw) validate(issuerID string) error {
	return validation.Errors{
		"path/{user-id}": validation.Validate(
			req.UserID, validation.Required, validation.By(MustBeValidID), validation.NotIn(issuerID),
		),
		"path/{claim-type}": validation.Validate(
			req.ClaimType, validation.Required, validation.By(MustBeClaimType),
		),
		"data/attributes/credential": validation.Validate(
			req.Body.Data.Attributes.CredentialSubject, validation.Required,
		),
		"data/attributes/expiration": validation.Validate(
			req.Body.Data.Attributes.Expiration,
			validation.When(
				!validation.IsEmpty(req.Body.Data.Attributes.Expiration),
				validation.By(MustBeValidRFC3339),
			),
		),
	}.Filter()
}

func MustBeClaimType(src interface{}) error {
	schemaTypeRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a schema type")
	}

	if _, ok := claims.ClaimSchemaTypeList[schemaTypeRaw]; !ok {
		return errors.New("schema type doesn't exist")
	}

	return nil
}

func MustBeValidRFC3339(src interface{}) error {
	expirationRAW, ok := src.(string)
	if !ok {
		return errors.New("it is not a valid an expiration")
	}

	if _, err := time.Parse(time.RFC3339, expirationRAW); err != nil {
		return errors.New("it is not a valid RFC3339 time format")
	}

	return nil
}

func (req *issueClaimRequestRaw) parse() *IssueClaimRequest {
	userID := core.ID{}

	_ = userID.UnmarshalText([]byte(req.UserID))
	did, _ := core.ParseDIDFromID(userID)

	schemaData, _ := req.Body.Data.Attributes.CredentialSubject.MarshalJSON()
	schemaDataTrimmed, _ := jsonRawTrimSpaces(schemaData)

	var expiration *time.Time
	if req.Body.Data.Attributes.Expiration != "" {
		parsedExpiration, _ := time.Parse(time.RFC3339, req.Body.Data.Attributes.Expiration)
		expiration = &parsedExpiration
	}

	return &IssueClaimRequest{
		Expiration: expiration,
		UserDID:    did,
		ClaimType:  claims.ClaimSchemaTypeList[req.ClaimType],
		Credential: schemaDataTrimmed,
	}
}

func jsonRawTrimSpaces(jsonRaw []byte) ([]byte, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(jsonRaw, &jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	result, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json")
	}

	return result, nil
}
