package pg

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"

	"gitlab.com/rarimo/identity/issuer/internal/data"
)

const (
	keyColumnName   = "key"
	valueColumnName = "value"
)

type treeStorageQ struct {
	db       *pgdb.DB
	treeName string
}

func NewTreeStorageQ(db *pgdb.DB, treeName string) data.TreeStorageQ {
	return &treeStorageQ{
		db:       db,
		treeName: treeName,
	}
}

func (q *treeStorageQ) New() data.TreeStorageQ {
	return NewTreeStorageQ(q.db.Clone(), q.treeName)
}

func (q *treeStorageQ) Insert(key, value []byte) error {
	err := q.db.Exec(
		sq.Insert(q.treeName).
			SetMap(map[string]interface{}{keyColumnName: key, valueColumnName: value}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *treeStorageQ) Upsert(key, value []byte) error {
	err := q.db.Exec(
		sq.Insert(q.treeName).
			SetMap(map[string]interface{}{keyColumnName: key, valueColumnName: value}).
			Suffix(
				fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s",
					keyColumnName, valueColumnName, valueColumnName),
			),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *treeStorageQ) Get(key []byte) ([]byte, error) {
	var result data.KeyValueStorage

	err := q.db.
		Get(&result, sq.Select("*").
			From(q.treeName).
			Where(sq.Eq{keyColumnName: key}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select rows")
	}

	return result.Value, nil
}
