package data

import (
	"database/sql/driver"

	core "github.com/iden3/go-iden3-core"
	"github.com/pkg/errors"
)

type ClaimsQ interface {
	New() ClaimsQ

	Get(id uint64) (*Claim, error)
	GetAuthClaim() (*Claim, error)
	GetBySchemaType(schemaType string, userID string) (*Claim, error)
	Insert(*Claim) error
	Update(*Claim) error
}

type Claim struct {
	ID             string     `db:"id"                structs:"id"`
	SchemaType     string     `db:"schema_type"       structs:"schema_type"`
	Revoked        bool       `db:"revoked"           structs:"revoked"`
	Credential     []byte     `db:"data"              structs:"data"`
	CoreClaim      *CoreClaim `db:"core_claim"        structs:"-"`
	MTP            []byte     `db:"-"                 structs:"-"`
	SignatureProof []byte     `db:"signature_proof"   structs:"signature_proof"`
	UserID         string     `db:"user_id"           structs:"user_id"`
}

type CoreClaim struct {
	*core.Claim
}

func NewCoreClaim(claim *core.Claim) *CoreClaim {
	return &CoreClaim{claim}
}

func (c *CoreClaim) Value() (driver.Value, error) {
	result, err := c.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal core claim to binary")
	}

	return result, nil
}

func (c *CoreClaim) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion src.([]byte) failed")
	}

	parsed := CoreClaim{
		Claim: &core.Claim{},
	}
	err := parsed.UnmarshalBinary(source)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal core claim from binary")
	}

	*c = parsed
	return nil
}
