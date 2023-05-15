package validation

import (
	"encoding/json"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type daoMembership struct {
	IsMember string `json:"is_member"`
}

type daoMembershipInt struct {
	IsMember int `json:"is_member"`
}

// nolint
func MustBeDAOMembership(schemaData interface{}) error {
	rawData, ok := schemaData.(json.RawMessage)
	if !ok {
		return errors.New("it is not a valid claim data")
	}

	var data daoMembership
	if err := json.Unmarshal(rawData, &data); err != nil {
		return errors.New("it is not a valid DAO membership credentials")
	}

	return validation.Errors{
		"data/attributes/credential/is_member": validation.Validate(
			data.IsMember, validation.Required, validation.By(MustBeBooleanInt),
		),
	}.Filter()
}

func ParseDAOMembership(rawData []byte) ([]byte, error) {
	var data daoMembership
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal DAO membership data")
	}

	isMember, err := strconv.ParseInt(data.IsMember, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse is_member field")
	}

	convertedData, err := json.Marshal(daoMembershipInt{
		IsMember: int(isMember),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DAO membership")
	}

	return convertedData, nil
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
