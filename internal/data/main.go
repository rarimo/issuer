package data

type MasterQ interface {
	New() MasterQ

	ClaimsQ() ClaimsQ
	CommittedStatesQ() CommittedStatesQ
	ClaimsOffersQ() ClaimsOffersQ

	Transaction(func() error) error
}
