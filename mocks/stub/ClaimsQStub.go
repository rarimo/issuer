package stub

import "gitlab.com/q-dev/q-id/issuer/internal/data"

type ClaimsQStub struct {
	NewStub             func() data.ClaimsQ
	GetStub             func(id string) (*data.Claim, error)
	GetAuthClaimStub    func() (*data.Claim, error)
	GetBySchemaTypeStub func(schemaType string, userID string) (*data.Claim, error)
	InsertStub          func(*data.Claim) error
	UpdateStub          func(*data.Claim) error
}

func (m *ClaimsQStub) New() data.ClaimsQ {
	return m.NewStub()
}

func (m *ClaimsQStub) Get(id string) (*data.Claim, error) {
	return m.GetStub(id)
}

func (m *ClaimsQStub) GetAuthClaim() (*data.Claim, error) {
	return m.GetAuthClaimStub()
}

func (m *ClaimsQStub) GetBySchemaType(schemaType string, userID string) (*data.Claim, error) {
	return m.GetBySchemaTypeStub(schemaType, userID)
}

func (m *ClaimsQStub) Insert(claim *data.Claim) error {
	return m.InsertStub(claim)
}

func (m *ClaimsQStub) Update(claim *data.Claim) error {
	return m.UpdateStub(claim)
}
