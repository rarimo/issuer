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
)

var (
	ErrSchemaFormatIsNotSupported = errors.New("schema format is not supported")
	ErrValidationData             = errors.New("data is not valid for requested schema")
)

type Builder struct {
	ipfsURL *url.URL
}
