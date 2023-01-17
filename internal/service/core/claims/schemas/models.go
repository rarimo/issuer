package schemas

import (
	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	httpProtocolName  = "http"
	httpsProtocolName = "https"
	ipfsProtocolName  = "ipfs"

	SchemaFormatJSONLD = "json-ld"
	SchemaFormatJSON   = "json"

	AuthBJJCredentialHash          = "ca938857241db9451ea329256b9c06e5"                                                                       //nolint
	AuthBJJCredentialSchemaURL     = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld"           //nolint
	KYCFullNameCredentialSchemaURL = "https://raw.githubusercontent.com/OmegaTymbJIep/schemas/main/kyc_credentials/kyc_credentials.json-ld"   //nolint
	QDAOMembershipSchemaURL        = "https://raw.githubusercontent.com/OmegaTymbJIep/schemas/main/q_dao_membership/q_dao_membership.json-ld" //nolint
)

var (
	ErrSchemaFormatIsNotSupported = errors.New("schema format is not supported")
	ErrValidationData             = errors.New("data is not valid for requested schema")
)

type Builder struct {
	ipfsURL *url.URL
}
