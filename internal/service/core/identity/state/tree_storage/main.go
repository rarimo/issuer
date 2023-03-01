package treestorage

import (
	"context"

	"github.com/iden3/go-merkletree-sql/v2"
	"gitlab.com/distributed_lab/kit/pgdb"
	errPkg "gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/data/pg"
)

// treeStorage implement merkletree.Storage
type treeStorage struct {
	currentRoot *merkletree.Hash
	db          data.TreeStorageQ
}

func NewTreeStorage(db *pgdb.DB, treeName string) merkletree.Storage {
	return &treeStorage{
		db: pg.NewTreeStorageQ(db, treeName),
	}
}

func (t *treeStorage) Get(_ context.Context, key []byte) (*merkletree.Node, error) {
	nodeRaw, err := t.db.Get(key)
	if err != nil {
		return nil, errPkg.Wrap(err, "failed to get node from the db")
	}

	node, err := merkletree.NewNodeFromBytes(nodeRaw)
	if err != nil {
		return nil, errPkg.Wrap(err, "failed to parse merkleTree node from bytes")
	}

	return node, nil
}

func (t *treeStorage) Put(_ context.Context, key []byte, node *merkletree.Node) error {
	if err := t.db.Insert(key, node.Value()); err != nil {
		return errPkg.Wrap(err, "failed to insert node into the db")
	}

	return nil
}

func (t *treeStorage) GetRoot(_ context.Context) (*merkletree.Hash, error) {
	var rootHash *merkletree.Hash
	if t.currentRoot != nil {
		copy(rootHash[:], t.currentRoot[:])

		return rootHash, nil
	}

	rootHashRaw, err := t.db.Get([]byte(rootKey))
	if err != nil {
		return nil, errPkg.Wrap(err, "failed to get the root merkle tree hash from the db")
	}
	if len(rootHashRaw) == 0 || rootHashRaw == nil {
		return nil, merkletree.ErrNotFound
	}

	rootHash, err = merkletree.NewHashFromHex(string(rootHashRaw))
	if err != nil {
		return nil, errPkg.Wrap(err, "failed to parse rootHash from hex")
	}

	t.currentRoot = &merkletree.Hash{}
	copy(t.currentRoot[:], rootHash[:])

	return rootHash, nil
}

func (t *treeStorage) SetRoot(_ context.Context, rootHash *merkletree.Hash) error {
	err := t.db.Upsert([]byte(rootKey), []byte(rootHash.Hex()))
	if err != nil {
		return errPkg.Wrap(err, "failed to insert new root hash")
	}

	return nil
}
