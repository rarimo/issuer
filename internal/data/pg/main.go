package pg

import (
	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/rarimo/identity/issuer/internal/data"
)

var sortColumns = map[string]string{
	createdAtColumnName: createdAtColumnName,
}

type masterQ struct {
	db *pgdb.DB
}

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{
		db: db,
	}
}

func (q *masterQ) New() data.MasterQ {
	return NewMasterQ(q.db.Clone())
}

func (q *masterQ) ClaimsQ() data.ClaimsQ {
	return NewClaimsQ(q.db)
}

func (q *masterQ) CommittedStatesQ() data.CommittedStatesQ {
	return NewCommittedStateQ(q.db)
}

func (q *masterQ) ClaimsOffersQ() data.ClaimsOffersQ {
	return NewClaimsOffersQ(q.db)
}

func (q *masterQ) Transaction(fn func() error) error {
	return q.db.Transaction(fn)
}
