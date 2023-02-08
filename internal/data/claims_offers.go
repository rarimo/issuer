package data

import "time"

type ClaimsOffersQ interface {
	New() ClaimsOffersQ

	Get(id string) (*ClaimOffer, error)
	Insert(claimOffer *ClaimOffer) error
	Update(claimOffer *ClaimOffer) error
}

type ClaimOffer struct {
	ID         string    `db:"id"          structs:"id"`
	From       string    `db:"from_id"     structs:"from_id"`
	To         string    `db:"to_id"       structs:"to_id"`
	CreatedAt  time.Time `db:"created_at"  structs:"created_at"`
	ClaimID    string    `db:"claim_id"    structs:"claim_id"`
	IsReceived bool      `db:"is_received" structs:"is_received"`
}
