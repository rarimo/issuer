package issuer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	resources "gitlab.com/q-dev/q-id/resources/claim_resources"
)

func (isr *issuer) compactClaim(
	ctx context.Context,
	userDID *core.DID,
	expiration *time.Time,
	claimType resources.ClaimSchemaType,
	credentialSubjectRaw []byte,
) (*data.Claim, error) {
	revNonce, err := claims.CryptoRandUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random uint64")
	}

	credentialsStatus := verifiable.CredentialStatus{
		ID:   fmt.Sprint(isr.baseURL, claims.CredentialStatusCheckURL, revNonce),
		Type: verifiable.SparseMerkleTreeProof,
	}

	claimID := uuid.NewString()
	credential, err := isr.newW3CCredential(claimID, userDID, credentialSubjectRaw, expiration, claimType, credentialsStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new w3c credential")
	}

	credentialRaw, err := json.Marshal(credential)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal credential")
	}

	coreClaim, err := isr.schemaBuilder.CreateCoreClaim(
		ctx,
		claimType,
		credential,
		revNonce,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process schema")
	}

	signProof, err := isr.signClaim(
		coreClaim,
		fmt.Sprint(isr.baseURL, claims.CredentialStatusCheckURL, coreClaim.GetRevocationNonce()),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signature proof")
	}

	return &data.Claim{
		ID:             claimID,
		SchemaType:     claimType.ToRaw(),
		CoreClaim:      data.NewCoreClaim(coreClaim),
		SignatureProof: signProof,
		Credential:     credentialRaw,
		UserID:         userDID.ID.String(),
	}, nil
}

func (isr *issuer) newW3CCredential(
	claimID string,
	userDID *core.DID,
	credentialSubjectRaw []byte,
	expiration *time.Time,
	claimType resources.ClaimSchemaType,
	credentialStatus verifiable.CredentialStatus,
) (*verifiable.W3CCredential, error) {
	issuanceDate := time.Now()
	credential, err := claims.ParseCredentialFromSnakeCase(credentialSubjectRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse credentials")
	}

	credential.CredentialSubject["id"] = userDID.String()
	credential.CredentialSubject["type"] = claimType.ToRaw()
	credential.ID = fmt.Sprint(isr.baseURL, GetClaimPath, claimID)
	credential.Expiration = expiration
	credential.IssuanceDate = &issuanceDate
	credential.Issuer = isr.Identifier.String()
	credential.CredentialStatus = credentialStatus
	credential.CredentialSchema = verifiable.CredentialSchema{
		ID:   resources.ClaimSchemaList[claimType].ClaimSchemaURL,
		Type: verifiable.JSONSchemaValidator2018,
	}
	credential.Context = []string{
		verifiable.JSONLDSchemaW3CCredential2018,
		verifiable.JSONLDSchemaIden3Credential,
		isr.schemaBuilder.CachedSchemas[string(claimType)].JSONLdContext,
	}
	credential.Type = []string{verifiable.TypeW3CVerifiableCredential, claimType.ToRaw()}

	return credential, nil
}

func (isr *issuer) signClaim(claim *core.Claim, checkRevLink string) ([]byte, error) {
	claimSign, err := claims.SignClaimEntry(claim, isr.Identity.Sign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign core claim")
	}

	signProof, err := claims.ConstructSignProof(isr.AuthClaim, claim, claimSign)
	if err != nil {
		return nil, errors.Wrap(err, "failed to construct sign proof")
	}
	signProof.IssuerData.CredentialStatus = &verifiable.CredentialStatus{
		ID:              checkRevLink,
		Type:            verifiable.SparseMerkleTreeProof,
		RevocationNonce: claim.GetRevocationNonce(),
	}

	signProofRaw, err := json.Marshal(signProof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal signature proof")
	}

	return signProofRaw, nil
}
