package claims

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
)

func ParseCoreClaim(data *CoreClaimData) (*core.Claim, error) {
	schemaBytes, err := hex.DecodeString(data.SchemaHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode schema hash from hex to bytes")
	}

	if len(schemaBytes) < CorrectSchemaHashBytesLength {
		return nil, ErrIncorrectSchemaHashLength
	}

	revNonce, err := CryptoRandUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random uint64")
	}

	var schemaHash core.SchemaHash
	copy(schemaHash[:], schemaBytes)
	coreClaim, err := core.NewClaim(
		schemaHash,
		core.WithIndexDataBytes(data.Slots.IndexA, data.Slots.IndexB),
		core.WithValueDataBytes(data.Slots.ValueA, data.Slots.ValueB),
		core.WithRevocationNonce(revNonce),
		core.WithExpirationDate(data.Expiration),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new claim")
	}

	if data.SubjectID != "" {
		userID, err := core.IDFromString(data.SubjectID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse user ID from string")
		}

		switch data.SubjectPosition {
		case "", SubjectPositionIndex:
			coreClaim.SetIndexID(userID)
		case SubjectPositionValue:
			coreClaim.SetValueID(userID)
		default:
			return nil, ErrIncorrectSubjectPosition
		}
	}

	return coreClaim, nil
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

func ConstructSignProof(authClaim *data.Claim, signature string) (*verifiable.BJJSignatureProof2021, error) {
	authMTP := &verifiable.Iden3SparseMerkleProof{}
	err := json.Unmarshal(authClaim.MTP, authMTP)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal auth core claim merkle tree proof")
	}

	authMTP.IssuerData.AuthClaim = authClaim.CoreClaim.Claim
	return &verifiable.BJJSignatureProof2021{
		Type:       BabyJubSignatureType,
		Signature:  signature,
		IssuerData: authMTP.IssuerData,
	}, nil
}

func GetCheckClaimRevLink(
	domain string,
	revNonce uint64,
	statusType verifiable.CredentialStatusType,
) ([]byte, string, error) {
	revLink := fmt.Sprint(domain, credentialStatusCheckURL, revNonce)
	credentialsStatus := verifiable.CredentialStatus{
		ID:   fmt.Sprint(domain, credentialStatusCheckURL, revNonce),
		Type: statusType,
	}

	credentialsStatusRaw, err := json.Marshal(credentialsStatus)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to marshal credential status")
	}

	return credentialsStatusRaw, revLink, nil
}
