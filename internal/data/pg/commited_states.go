package pg

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/rarimo/issuer/internal/data"
)

const (
	committedStatesTableName = "committed_states"
	txIDColumnName           = "tx_id"
	createdAtColumnName      = "created_at"
	statusColumnName         = "status"
	isGenesisColumnName      = "is_genesis"

	descCreatedAtColumnName = "-" + createdAtColumnName
)

type committedStatesQ struct {
	db  *pgdb.DB
	sel sq.SelectBuilder
}

func NewCommittedStateQ(db *pgdb.DB) data.CommittedStatesQ {
	return &committedStatesQ{
		db:  db,
		sel: sq.Select("*").From(committedStatesTableName),
	}
}

func (q *committedStatesQ) New() data.CommittedStatesQ {
	return NewCommittedStateQ(q.db.Clone())
}

func (q *committedStatesQ) Select() ([]data.CommittedState, error) {
	var result []data.CommittedState
	err := q.db.Select(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return result, nil
}

func (q *committedStatesQ) Get(id uint64) (*data.CommittedState, error) {
	var result data.CommittedState

	err := q.db.Get(&result,
		sq.Select("*").
			From(committedStatesTableName).
			Where(sq.Eq{idColumnName: id}),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}

func (q *committedStatesQ) Sort(sort pgdb.SortedOffsetPageParams) data.CommittedStatesQ {
	q.sel = sort.ApplyTo(q.sel, sortColumns)
	return q
}

func (q *committedStatesQ) GetLatest() (*data.CommittedState, error) {
	result, err := q.Sort(pgdb.SortedOffsetPageParams{
		Limit: 1,
		Sort:  []string{descCreatedAtColumnName},
	}).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select rows")
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func (q *committedStatesQ) GetGenesis() (*data.CommittedState, error) {
	var result data.CommittedState

	err := q.db.Get(&result,
		sq.Select("*").
			From(committedStatesTableName).
			Where(sq.Eq{isGenesisColumnName: true}),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}

func (q *committedStatesQ) Insert(committedState *data.CommittedState) error {
	err := q.db.Get(&committedState.ID,
		sq.Insert(committedStatesTableName).
			SetMap(structs.Map(committedState)).
			Suffix(fmt.Sprintf("returning %s", idColumnName)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *committedStatesQ) Update(committedState *data.CommittedState) error {
	err := q.db.Exec(
		sq.Update(committedStatesTableName).
			SetMap(structs.Map(committedState)).
			Where(sq.Eq{idColumnName: committedState.ID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update rows")
	}

	return nil
}

func (q *committedStatesQ) WhereStatus(status data.Status) data.CommittedStatesQ {
	q.sel = q.sel.Where(sq.Eq{statusColumnName: status})
	return q
}
