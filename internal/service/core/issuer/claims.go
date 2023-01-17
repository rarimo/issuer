package issuer

import (
	"context"
	"encoding/json"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
)

func (isr *issuer) compactClaim(
	ctx context.Context,
	userID *core.ID,
	expiration time.Time,
	schemaType claims.ClaimSchemaType,
	schemaData []byte,
) (*data.Claim, error) {
	slots, schemaHash, err := isr.schemaBuilder.Process(
		ctx,
		schemaData,
		schemaType.ToRaw(),
		claims.ClaimSchemaList[schemaType].ClaimSchemaURL,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process schema")
	}

	coreClaim, err := claims.ParseCoreClaim(&claims.CoreClaimData{
		SchemaHash:      schemaHash,
		Slots:           *slots,
		SubjectID:       userID.String(),
		Expiration:      expiration,
		SubjectPosition: claims.SubjectPositionIndex,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse core claim")
	}

	credentialsStatus, checkRevLink, err := claims.GetCheckClaimRevLink(
		isr.domain,
		coreClaim.GetRevocationNonce(),
		verifiable.SparseMerkleTreeProof,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create credential status")
	}

	signProof, err := isr.signClaim(coreClaim, checkRevLink)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signature proof")
	}

	return &data.Claim{
		SchemaURL:        claims.ClaimSchemaList[schemaType].ClaimSchemaURL,
		SchemaType:       schemaType.ToRaw(),
		CoreClaim:        data.NewCoreClaim(coreClaim),
		CredentialStatus: credentialsStatus,
		SignatureProof:   signProof,
		Data:             schemaData,
	}, nil
}

func (isr *issuer) signClaim(claim *core.Claim, checkRevLink string) ([]byte, error) {
	claimSign, err := claims.SignClaimEntry(claim, isr.Identity.Sign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign core claim")
	}

	signProof, err := claims.ConstructSignProof(isr.AuthClaim, claimSign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct sign proof")
	}
	signProof.IssuerData.RevocationStatus = checkRevLink

	signProofRaw, err := json.Marshal(signProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature proof")
	}

	return signProofRaw, nil
}
