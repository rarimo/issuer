package claims

import (
	"encoding/json"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewAuthClaim(key *babyjub.PublicKey, schemaHash core.SchemaHash) (*core.Claim, error) {
	revNonce, err := CryptoRandUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random uint64")
	}

	claim, err := core.NewClaim(
		schemaHash,
		core.WithIndexDataInts(key.X, key.Y),
		core.WithRevocationNonce(revNonce),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new core claim")
	}

	return claim, nil
}

func GenerateAuthClaimData(publicKey *babyjub.PublicKey) ([]byte, error) {
	jsonRaw, err := json.Marshal(&map[string]string{
		"x": publicKey.X.String(),
		"y": publicKey.Y.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal json structure")
	}

	return jsonRaw, nil
}
