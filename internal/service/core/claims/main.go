package claims

import (
	"encoding/hex"
	"encoding/json"
	"math/big"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/poseidon"
	jsonSuite "github.com/iden3/go-schema-processor/json"
	"github.com/iden3/go-schema-processor/utils"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

func DefineMerklizedRootPosition(metadata *jsonSuite.SchemaMetadata, position string) string {
	if metadata != nil && metadata.Serialization != nil {
		return ""
	}

	if position != "" {
		return position
	}

	return utils.MerklizedRootPositionIndex
}

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

func ConstructSignProof(authClaim *data.Claim, claim *core.Claim, signature string) (*verifiable.BJJSignatureProof2021, error) {
	authMTP := &verifiable.Iden3SparseMerkleProof{}
	err := json.Unmarshal(authClaim.MTP, authMTP)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal auth core claim merkle tree proof")
	}

	authCoreClaimHex, err := authClaim.CoreClaim.Claim.Hex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hex from auth core claim")
	}

	coreClaimHex, err := claim.Hex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hex from auth core claim")
	}

	authMTP.IssuerData.AuthCoreClaim = authCoreClaimHex
	authMTP.IssuerData.MTP = authMTP.MTP
	return &verifiable.BJJSignatureProof2021{
		Type:       BabyJubSignatureType,
		Signature:  signature,
		CoreClaim:  coreClaimHex,
		IssuerData: authMTP.IssuerData,
	}, nil
}
