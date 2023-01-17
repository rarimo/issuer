package requests

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	core "github.com/iden3/go-iden3-core"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/resources"
)

type IssueClaimRequest struct {
	UserID     *core.ID
	SchemaType claims.ClaimSchemaType
	Expiration time.Time
	SchemaData []byte
}

type issueClaimRequestRaw struct {
	UserID string
	Body   resources.IssueClaimRequest
}

func NewIssueClaim(r *http.Request) (*IssueClaimRequest, error) {
	requestBody := resources.IssueClaimRequest{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, errors.Wrap(err, "failed to decode json request body")
	}

	requestRaw := issueClaimRequestRaw{
		UserID: chi.URLParam(r, UserIDPathParam),
		Body:   requestBody,
	}

	if err := requestRaw.validate(); err != nil {
		return nil, err
	}

	schemaType := requestRaw.Body.Data.Attributes.SchemaType
	if err := validation.Validate(
		requestRaw.Body.Data.Attributes.SchemaData,
		validation.By(claims.ClaimSchemaList[claims.ClaimSchemaTypeList[schemaType]].ClaimDataValidateFunc),
	); err != nil {
		return nil, errors.Wrap(err, "invalid schema data")
	}

	parseData, err := claims.ClaimSchemaList[claims.ClaimSchemaTypeList[schemaType]].ClaimDataParseFunc(
		requestRaw.Body.Data.Attributes.SchemaData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse kyc full name data")
	}

	requestRaw.Body.Data.Attributes.SchemaData = parseData

	return requestRaw.parse(), nil
}

// nolint
func (req *issueClaimRequestRaw) validate() error {
	return validation.Errors{
		"path/{user-id}": validation.Validate(
			req.UserID, validation.Required, validation.By(MustBeIden3Identifier),
		),
		"data/attributes/schema_type": validation.Validate(
			req.Body.Data.Attributes.SchemaType, validation.Required, validation.By(MustBeSchemaType),
		),
		"data/attributes/schema_data": validation.Validate(
			req.Body.Data.Attributes.SchemaData, validation.Required,
		),
		"data/attributes/expiration": validation.Validate(
			req.Body.Data.Attributes.Expiration, validation.Required, validation.By(MustBeValidRFC3339),
		),
	}.Filter()
}

func MustBeSchemaType(src interface{}) error {
	schemaTypeRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a schema type")
	}

	if _, ok := claims.ClaimSchemaTypeList[schemaTypeRaw]; !ok {
		return errors.New("it is not a schema type")
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
	userID := &core.ID{}

	_ = userID.UnmarshalText([]byte(req.UserID))

	schemaData, _ := req.Body.Data.Attributes.SchemaData.MarshalJSON()
	schemaDataTrimmed, _ := jsonRawTrimSpaces(schemaData)

	expiration, _ := time.Parse(time.RFC3339, req.Body.Data.Attributes.Expiration)

	return &IssueClaimRequest{
		Expiration: expiration,
		UserID:     userID,
		SchemaType: claims.ClaimSchemaTypeList[req.Body.Data.Attributes.SchemaType],
		SchemaData: schemaDataTrimmed,
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