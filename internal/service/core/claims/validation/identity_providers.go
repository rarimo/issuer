package validation

import (
	"encoding/json"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type identityProviders struct {
	Provider                 string `json:"provider"`
	Address                  string `json:"address"`
	GitcoinPassportScore     string `json:"gitcoin_passport_score"`
	WorldCoinScore           string `json:"worldcoin_score"`
	UnstoppableDomain        string `json:"unstoppable_domain"`
	CivicGatekeeperNetworkID string `json:"civic_gatekeeper_network_id"`
	KYCAdditionalData        string `json:"kyc_additional_data"`
	IsNaturalPerson          string `json:"is_natural"`
}

type identityProvidersParsed struct {
	Provider                 string `json:"provider"`
	Address                  string `json:"address"`
	GitcoinPassportScore     string `json:"gitcoin_passport_score"`
	WorldCoinScore           string `json:"worldcoin_score"`
	UnstoppableDomain        string `json:"unstoppable_domain"`
	CivicGatekeeperNetworkID int    `json:"civic_gatekeeper_network_id"`
	KYCAdditionalData        string `json:"kyc_additional_data"`
	IsNaturalPerson          int    `json:"is_natural"`
}

func MustBeIdentityProvidersCredentials(credentialSubject interface{}) error {
	rawData, ok := credentialSubject.(json.RawMessage)
	if !ok {
		return errors.New("it is not a valid credential subject")
	}

	var data identityProviders
	if err := json.Unmarshal(rawData, &data); err != nil {
		return errors.New("it is not a valid Identity Providers credentials")
	}

	return validation.Errors{
		"data/attributes/credential/provider": validation.Validate(
			data.Provider, validation.Required,
		),
		"data/attributes/credential/civic_gatekeeper_network_id": validation.Validate(
			data.CivicGatekeeperNetworkID, validation.By(MustBeUintOrEmpty),
		),
		"data/attributes/credential/is_natural": validation.Validate(
			data.IsNaturalPerson, validation.Required, validation.By(MustBeBooleanInt),
		),
	}.Filter()
}

func ParseIdentityProvidersCredentials(rawData []byte) ([]byte, error) {
	var data identityProviders
	err := json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal DAO membership data")
	}

	civicGatekeeperNetworkID, err := parseIntFromString(data.CivicGatekeeperNetworkID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse civic_gatekeeper_network_id field")
	}

	gitcoinPassportScore := data.GitcoinPassportScore
	if data.GitcoinPassportScore == "" {
		gitcoinPassportScore = "0.0"
	}

	isNaturalPerson, err := strconv.ParseInt(data.IsNaturalPerson, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse is_member field")
	}

	parsedCredentials, err := json.Marshal(identityProvidersParsed{
		Provider:                 noneIfEmpty(data.Provider),
		Address:                  noneIfEmpty(data.Address),
		WorldCoinScore:           noneIfEmpty(data.WorldCoinScore),
		UnstoppableDomain:        noneIfEmpty(data.UnstoppableDomain),
		GitcoinPassportScore:     gitcoinPassportScore,
		CivicGatekeeperNetworkID: civicGatekeeperNetworkID,
		KYCAdditionalData:        noneIfEmpty(data.KYCAdditionalData),
		IsNaturalPerson:          int(isNaturalPerson),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal DAO membership")
	}

	return parsedCredentials, nil
}

func parseIntFromString(src string) (int, error) {
	if src == "" {
		return 0, nil
	}

	gitcoinPassportScore, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse gitcoin_passport_score field")
	}

	return int(gitcoinPassportScore), nil
}

func MustBeUintOrEmpty(src interface{}) error {
	numberRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	if numberRaw == "" {
		return nil
	}

	_, err := strconv.ParseInt(numberRaw, 10, 64)
	if err != nil {
		return errors.New("it is not an uint64")
	}

	return nil
}

func noneIfEmpty(src string) string {
	if src == "" {
		return "none"
	}

	return src
}
