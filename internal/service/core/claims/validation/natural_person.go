package validation

import (
	"encoding/json"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type naturalPerson struct {
	IsNaturalPerson string `json:"is_natural_person"`
}

type naturalPersonParsed struct {
	IsNaturalPerson int `json:"is_natural_person"`
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

	return validation.Errors{
		"data/attributes/credential/is_natural_person": validation.Validate(
			data.IsNaturalPerson, validation.Required, validation.By(MustBeBooleanInt),
		),
	}.Filter()
}

func ParseNaturalPersonCredentials(rawData []byte) ([]byte, error) {
	var data naturalPerson
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal DAO membership data")
	}

	isNaturalPerson, err := strconv.ParseInt(data.IsNaturalPerson, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse is_member field")
	}

	parsedCredentials, err := json.Marshal(naturalPersonParsed{
		IsNaturalPerson: int(isNaturalPerson),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DAO membership")
	}

	return parsedCredentials, nil
}

func MustBeBooleanInt(src interface{}) error {
	numberRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	booleanInt, err := strconv.ParseInt(numberRaw, 10, 64)
	if err != nil {
		return errors.New("it is not an int64")
	}

	if booleanInt != 0 && booleanInt != 1 {
		return errors.New("it is not a boolean in integer format")
	}

	return nil
}
