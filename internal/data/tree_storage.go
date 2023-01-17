package data

type TreeStorageQ interface {
	New() TreeStorageQ

	Get(key []byte) ([]byte, error)
	Insert(key, value []byte) error
	Upsert(key, value []byte) error
}

type KeyValueStorage struct {
	Key   []byte `db:"key"`
	Value []byte `db:"value"`
}
