package data

import (
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type CommittedStatesQ interface {
	New() CommittedStatesQ

	Get(id uint64) (*CommittedState, error)
	Insert(committedState *CommittedState) error
	Select() ([]CommittedState, error)
	Update(committedState *CommittedState) error
	Sort(sort pgdb.SortedOffsetPageParams) CommittedStatesQ

	GetLatest() (*CommittedState, error)
	GetGenesis() (*CommittedState, error)
	WhereStatus(status Status) CommittedStatesQ
}

type CommittedState struct {
	ID                  uint64    `db:"id"                    structs:"-"`
	Status              Status    `db:"status"                structs:"status"`
	Message             string    `db:"message"               structs:"message"`
	TxID                string    `db:"tx_id"                 structs:"tx_id"`
	CreatedAt           time.Time `db:"created_at"            structs:"created_at"`
	BlockTimestamp      uint64    `db:"block_timestamp"       structs:"block_timestamp"`
	BlockNumber         uint64    `db:"block_number"          structs:"block_number"`
	IsGenesis           bool      `db:"is_genesis"            structs:"is_genesis"`
	RootsTreeRoot       []byte    `db:"roots_tree_root"       structs:"roots_tree_root"`
	ClaimsTreeRoot      []byte    `db:"claims_tree_root"      structs:"claims_tree_root"`
	RevocationsTreeRoot []byte    `db:"revocations_tree_root" structs:"revocations_tree_root"`
}

type Status string

const (
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusPending    = "pending"
	StatusFailed     = "failed"
)
