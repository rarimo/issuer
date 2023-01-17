package validation

import (
	"encoding/json"
	"math/big"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	NameRegExpr      = "^[A-Z]{1}[a-z]{1,15}$"
	NameAffixRegExpr = "^[A-Za-z]{1,8}$"
)

type kycFullNameCredentials struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	NameSuffix string `json:"name_suffix"`
	NamePrefix string `json:"name_prefix"`
}

type kycFullNameCredentialsInt struct {
	FirstName  *big.Int `json:"first_name"`
	LastName   *big.Int `json:"last_name"`
	MiddleName *big.Int `json:"middle_name"`
	NameSuffix *big.Int `json:"name_suffix"`
	NamePrefix *big.Int `json:"name_prefix"`
}

// nolint
func MustBeKYCFullNameCredentials(schemaData interface{}) error {
	data, ok := schemaData.(json.RawMessage)
	if !ok {
		return errors.New("it is not a valid claim data")
	}

	var dataRaw kycFullNameCredentials
	if err := json.Unmarshal(data, &dataRaw); err != nil {
		return errors.New("it is not a valid KYCNameCredentials")
	}

	return validation.Errors{
		"data/attributes/data/first_name": validation.Validate(
			dataRaw.FirstName, validation.Required, validation.Match(regexp.MustCompile(NameRegExpr)),
		),
		"data/attributes/data/last_name": validation.Validate(
			dataRaw.LastName, validation.Required, validation.Match(regexp.MustCompile(NameRegExpr)),
		),
		"data/attributes/data/middle_name": validation.Validate(
			dataRaw.MiddleName, validation.Match(regexp.MustCompile(NameRegExpr)),
		),
		"data/attributes/data/name_suffix": validation.Validate(
			dataRaw.NameSuffix, validation.Match(regexp.MustCompile(NameAffixRegExpr)),
		),
		"data/attributes/data/name_prefix": validation.Validate(
			dataRaw.NamePrefix, validation.Match(regexp.MustCompile(NameAffixRegExpr)),
		),
	}.Filter()
}

func ConvertKYCFullNameCredentials(fullNameStringDataRaw []byte) ([]byte, error) {
	var fullNameStringData kycFullNameCredentials
	err := json.Unmarshal(fullNameStringDataRaw, &fullNameStringData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal kyc full name credentials")
	}

	fullNameIntDataRaw, err := json.Marshal(kycFullNameCredentialsInt{
		FirstName:  stringToBigInt(fullNameStringData.FirstName),
		LastName:   stringToBigInt(fullNameStringData.LastName),
		MiddleName: stringToBigInt(fullNameStringData.MiddleName),
		NameSuffix: stringToBigInt(fullNameStringData.NameSuffix),
		NamePrefix: stringToBigInt(fullNameStringData.NamePrefix),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal kyc full name credentials")
	}

	return fullNameIntDataRaw, nil
}

func stringToBigInt(stringData string) *big.Int {
	return new(big.Int).SetBytes([]byte(stringData))
}
