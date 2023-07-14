package pg

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/rarimo/identity/issuer/internal/data"
)

const (
	claimsOffersTableName = "claims_offers"
)

type claimsOffersQ struct {
	db  *pgdb.DB
	sel sq.SelectBuilder
}

func NewClaimsOffersQ(db *pgdb.DB) data.ClaimsOffersQ {
	return &claimsOffersQ{
		db:  db,
		sel: sq.Select("*").From(claimsOffersTableName),
	}
}

func (q *claimsOffersQ) New() data.ClaimsOffersQ {
	return NewClaimsOffersQ(q.db.Clone())
}

func (q *claimsOffersQ) Get(id string) (*data.ClaimOffer, error) {
	var result data.ClaimOffer

	err := q.db.Get(&result,
		sq.Select("*").
			From(claimsOffersTableName).
			Where(sq.Eq{idColumnName: id}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}

func (q *claimsOffersQ) Insert(claimOffer *data.ClaimOffer) error {
	err := q.db.Exec(sq.Insert(claimsOffersTableName).SetMap(structs.Map(claimOffer)))
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *claimsOffersQ) Update(claimOffer *data.ClaimOffer) error {
	err := q.db.Exec(
		sq.Update(claimsOffersTableName).
			SetMap(structs.Map(claimOffer)).
			Where(sq.Eq{idColumnName: claimOffer.ID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}
