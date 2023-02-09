package stub

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

type CommittedStatesQStub struct {
	NewStub         func() data.CommittedStatesQ
	GetStub         func(id uint64) (*data.CommittedState, error)
	InsertStub      func(*data.CommittedState) error
	SelectStub      func() ([]data.CommittedState, error)
	UpdateStub      func(*data.CommittedState) error
	SortStub        func(sort pgdb.SortedOffsetPageParams) data.CommittedStatesQ
	GetLatestStub   func() (*data.CommittedState, error)
	GetGenesisStub  func() (*data.CommittedState, error)
	WhereStatusStub func(status data.Status) data.CommittedStatesQ
}

func (m *CommittedStatesQStub) New() data.CommittedStatesQ {
	return m.NewStub()
}

func (m *CommittedStatesQStub) Get(id uint64) (*data.CommittedState, error) {
	return m.GetStub(id)
}

func (m *CommittedStatesQStub) Insert(committedState *data.CommittedState) error {
	return m.InsertStub(committedState)
}

func (m *CommittedStatesQStub) Select() ([]data.CommittedState, error) {
	return m.SelectStub()
}

func (m *CommittedStatesQStub) Update(committedState *data.CommittedState) error {
	return m.UpdateStub(committedState)
}

func (m *CommittedStatesQStub) Sort(sort pgdb.SortedOffsetPageParams) data.CommittedStatesQ {
	return m.SortStub(sort)
}

func (m *CommittedStatesQStub) GetLatest() (*data.CommittedState, error) {
	return m.GetLatestStub()
}

func (m *CommittedStatesQStub) GetGenesis() (*data.CommittedState, error) {
	return m.GetGenesisStub()
}

func (m *CommittedStatesQStub) WhereStatus(status data.Status) data.CommittedStatesQ {
	return m.WhereStatusStub(status)
}
