package claims

import (
	"encoding/hex"
	"math/big"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/poseidon"
	jsonSuite "github.com/iden3/go-schema-processor/json"
	"github.com/iden3/go-schema-processor/utils"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func SignClaimEntry(claim *core.Claim, signFunc func(*big.Int) ([]byte, error)) (string, error) {
	hashIndex, hashValue, err := claim.HiHv()
	if err != nil {
		return "", errors.Wrap(err, "failed to get hash of the index and value from the claim")
	}

	commonHash, err := poseidon.Hash([]*big.Int{hashIndex, hashValue})
	if err != nil {
		return "", errors.Wrap(err, "failed to poseidon hash the index and value of the claim")
	}

	sig, err := signFunc(commonHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign the claim")
	}

	return hex.EncodeToString(sig), nil
}

func DefineMerklizedRootPosition(metadata *jsonSuite.SchemaMetadata, position string) string {
	if metadata != nil && metadata.Serialization != nil {
		return ""
	}

	if position != "" {
		return position
	}

	return utils.MerklizedRootPositionIndex
}
