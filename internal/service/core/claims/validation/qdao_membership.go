package validation

import (
	"encoding/json"
	"math/big"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	BoolNumberRegExpr = "^[10]{1}$"
)

type qDaoMembership struct {
	IsMember string `json:"is_member"`
}

type qDaoMembershipInt struct {
	IsMember *big.Int `json:"is_member"`
}

// nolint
func MustBeQDAOMembership(schemaData interface{}) error {
	rawData, ok := schemaData.(json.RawMessage)
	if !ok {
		return errors.New("it is not a valid claim data")
	}

	var data qDaoMembership
	if err := json.Unmarshal(rawData, &data); err != nil {
		return errors.New("it is not a valid Q DAO membership credentials")
	}

	return validation.Errors{
		"data/attributes/data/is_member": validation.Validate(
			data.IsMember, validation.Required, validation.Match(regexp.MustCompile(BoolNumberRegExpr)),
		),
	}.Filter()
}

func ConvertQDAOMembership(rawData []byte) ([]byte, error) {
	var data qDaoMembership
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Q DAO membership data")
	}

	isMember, ok := new(big.Int).SetString(data.IsMember, 10) //nolint
	if !ok {
		return nil, errors.New("failed to parse is_member field")
	}

	convertedData, err := json.Marshal(qDaoMembershipInt{
		IsMember: isMember,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Q DAO membership")
	}

	return convertedData, nil
}
