package requests

import (
	"encoding/json"
	"io"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/iden3/go-jwz"
	"github.com/iden3/iden3comm/protocol"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type OfferCallbackRequest struct {
	Token        *jwz.Token
	FetchMessage *protocol.CredentialFetchRequestMessage
}

func NewOfferCallback(r *http.Request) (*OfferCallbackRequest, error) {
	tokenRaw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("is not a valid request token")
	}

	token, err := jwz.Parse(string(tokenRaw))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse jwz")
	}

	var fetchMessage protocol.CredentialFetchRequestMessage
	if err := json.Unmarshal(token.GetPayload(), &fetchMessage); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	requestBody := OfferCallbackRequest{
		Token:        token,
		FetchMessage: &fetchMessage,
	}
	if err := requestBody.validate(); err != nil {
		return nil, err
	}

	return &requestBody, nil
}

// nolint
func (r *OfferCallbackRequest) validate() error {
	return validation.Errors{
		"message/type": validation.Validate(
			r.FetchMessage.Type, validation.Required, validation.In(protocol.CredentialFetchRequestMessageType),
		),
		"message/from": validation.Validate(
			r.FetchMessage.From, validation.Required, validation.By(MustBeValidDID),
		),
		"message/to": validation.Validate(
			r.FetchMessage.To, validation.Required, validation.By(MustBeValidDID),
		),
		"message/body/id": validation.Validate(
			r.FetchMessage.Body.ID, validation.Required, validation.By(MustBeClaimID),
		),
	}.Filter()
}

func MustBeClaimID(src interface{}) error {
	uuidRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	_, err := uuid.Parse(uuidRaw)
	if err != nil {
		return errors.New("it is not a valid uuid")
	}

	return nil
}
