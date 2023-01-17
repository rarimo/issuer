package pg

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

const (
	claimsTableName     = "claims"
	idColumnName        = "id"
	coreClaimColumnName = "core_claim"

	authClaimID = 1
)

type claimsQ struct {
	db *pgdb.DB
}

func NewClaimsQ(db *pgdb.DB) data.ClaimsQ {
	return &claimsQ{
		db: db.Clone(),
	}
}

func (q *claimsQ) New() data.ClaimsQ {
	return NewClaimsQ(q.db)
}

func (q *claimsQ) Insert(claim *data.Claim) error {
	clauses := structs.Map(claim)
	clauses[coreClaimColumnName] = claim.CoreClaim

	err := q.db.Get(&claim.ID,
		sq.Insert(claimsTableName).
			SetMap(clauses).
			Suffix(fmt.Sprintf("returning %s", idColumnName)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *claimsQ) Get(id uint64) (*data.Claim, error) {
	var result data.Claim

	err := q.db.Get(&result,
		sq.Select("*").
			From(claimsTableName).
			Where(sq.Eq{idColumnName: id}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}

func (q *claimsQ) GetAuthClaim() (*data.Claim, error) {
	var result data.Claim

	err := q.db.Get(&result,
		sq.Select("*").
			From(claimsTableName).
			Where(sq.Eq{idColumnName: authClaimID}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}