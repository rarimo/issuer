package schemas

import (
	"time"

	jsonSuite "github.com/iden3/go-schema-processor/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrValidationData = errors.New("data is not valid for requested schema")
)

type Builder struct {
	SchemasBaseURL string
	CachedSchemas  map[string]Schema
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
