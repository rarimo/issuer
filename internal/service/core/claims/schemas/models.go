package schemas

import (
	"time"

	jsonSuite "github.com/iden3/go-schema-processor/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	AuthBJJCredentialSchemaURL = "https://schema.iden3.io/core/jsonld/auth.jsonld#AuthBJJCredential" //nolint
)

var (
	ErrValidationData = errors.New("data is not valid for requested schema")
)

type Builder struct {
	CachedSchemas map[string]Schema
}

type Schema struct {
	Raw           []byte
	Body          jsonSuite.Schema
	JSONLdContext string
}

type CompactClaimOptions struct {
	Expiration *time.Time
	Version    uint32
}
