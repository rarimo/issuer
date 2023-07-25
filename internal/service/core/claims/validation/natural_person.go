package validation

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

type naturalPerson struct {
	IsNaturalPerson bool `json:"is_natural"`
}

type naturalPersonParsed struct {
	IsNaturalPerson bool `json:"is_natural"`
}

// nolint
func MustBeNaturalPersonCredentials(credentialSubject interface{}) error {
	rawData, ok := credentialSubject.(json.RawMessage)
	if !ok {
		return errors.New("it is not a valid credential subject")
	}

	var data naturalPerson
	if err := json.Unmarshal(rawData, &data); err != nil {
		return errors.New("it is not a valid Natural person credentials")
	}

	return nil
}

func ParseNaturalPersonCredentials(rawData []byte) ([]byte, error) {
	var data naturalPerson
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal DAO membership data")
	}

	parsedCredentials, err := json.Marshal(naturalPersonParsed{
		IsNaturalPerson: data.IsNaturalPerson,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DAO membership")
	}

	return parsedCredentials, nil
}
