package identity

import (
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/rarimo/identity/issuer/internal/data"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/identity/state"
)

var (
	ErrClaimWasNotPublishedYet = errors.New("claim was not published yet")
)

type Identity struct {
	babyJubJubPrivateKey *babyjub.PrivateKey
	Identifier           *core.DID
	AuthClaim            *data.Claim
	circuitsPath         string

	log   *logan.Entry
	State *state.IdentityState
}

type Iden3SparseMerkleTreeProofWithID struct {
	verifiable.Iden3SparseMerkleTreeProof

	ID string `json:"id"`
}
