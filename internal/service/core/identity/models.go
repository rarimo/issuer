package identity

import (
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
)

type Identity struct {
	babyJubJubPrivateKey *babyjub.PrivateKey
	Identifier           *core.ID
	AuthClaim            *data.Claim
	circuitsPath         string

	log   *logan.Entry
	State *state.IdentityState
}
