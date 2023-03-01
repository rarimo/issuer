package state

import (
	"sync"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/publisher"
)

const (
	StateTransitionCircuitWasmPath = "/state_transition/circuit.wasm"
	StateTransitionCircuitFinalKey = "/state_transition/circuit_final.zkey"
)

var (
	ErrStateWasntChanged = errors.New("state was not changed")
	ErrOldStateNotFound  = errors.New("old state not found")
)

type IdentityState struct {
	ClaimsQ         data.ClaimsQ
	CommittedStateQ data.CommittedStatesQ
	ClaimsTree      *merkletree.MerkleTree
	RevocationsTree *merkletree.MerkleTree
	RootsTree       *merkletree.MerkleTree

	circuits     map[string][]byte
	publisher    publisher.Publisher
	circuitsPath string

	*sync.Mutex
}

type Config struct {
	DB              *pgdb.DB
	PublisherConfig *publisher.Config
	IdentityConfig  *config.IdentityConfig
}

type CommittedState struct {
	ID      uint64
	Status  data.Status
	Message string

	CommitInfo *CommitInfo
	CreatedAt  time.Time

	IsGenesis           bool
	RootsTreeRoot       *merkletree.Hash
	ClaimsTreeRoot      *merkletree.Hash
	RevocationsTreeRoot *merkletree.Hash
}

type CommitInfo struct {
	TxID           string
	BlockTimestamp uint64
	BlockNumber    uint64
}

type IdentityInfo struct {
	BabyJubJubPrivateKey *babyjub.PrivateKey
	Identifier           *core.ID
	AuthClaim            *core.Claim
}

func (cs *CommittedState) ToRaw() *data.CommittedState {
	result := data.CommittedState{
		ID:        cs.ID,
		CreatedAt: cs.CreatedAt,
		IsGenesis: cs.IsGenesis,
		Status:    cs.Status,
		Message:   cs.Message,
	}

	if cs.CommitInfo != nil {
		result.BlockTimestamp = cs.CommitInfo.BlockTimestamp
		result.BlockNumber = cs.CommitInfo.BlockNumber
		result.TxID = cs.CommitInfo.TxID
	}

	if cs.RootsTreeRoot != nil {
		result.RootsTreeRoot = append([]byte{}, cs.RootsTreeRoot[:]...)
	}

	if cs.ClaimsTreeRoot != nil {
		result.ClaimsTreeRoot = append([]byte{}, cs.ClaimsTreeRoot[:]...)
	}

	if cs.RevocationsTreeRoot != nil {
		result.RevocationsTreeRoot = append([]byte{}, cs.RevocationsTreeRoot[:]...)
	}

	return &result
}

func CommittedStateFromRaw(rawState *data.CommittedState) (*CommittedState, error) {
	if rawState == nil {
		return nil, errors.New("rawState is nil")
	}

	result := &CommittedState{
		ID:        rawState.ID,
		Status:    rawState.Status,
		Message:   rawState.Message,
		CreatedAt: rawState.CreatedAt,
		IsGenesis: rawState.IsGenesis,

		CommitInfo: &CommitInfo{
			TxID:           rawState.TxID,
			BlockTimestamp: rawState.BlockTimestamp,
			BlockNumber:    rawState.BlockNumber,
		},

		ClaimsTreeRoot:      &merkletree.Hash{},
		RevocationsTreeRoot: &merkletree.Hash{},
		RootsTreeRoot:       &merkletree.Hash{},
	}

	copy(result.ClaimsTreeRoot[:], rawState.ClaimsTreeRoot)
	copy(result.RevocationsTreeRoot[:], rawState.RevocationsTreeRoot)
	copy(result.RootsTreeRoot[:], rawState.RootsTreeRoot)

	return result, nil
}
