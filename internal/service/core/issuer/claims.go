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
	userID *core.ID,
	expiration *time.Time,
	schemaType resources.ClaimSchemaType,
	credentialRaw []byte,
) (*data.Claim, error) {
	revNonce, err := claims.CryptoRandUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random uint64")
	}

	credentialsStatus := verifiable.CredentialStatus{
		ID:   fmt.Sprint(isr.baseURL, claims.CredentialStatusCheckURL, revNonce),
		Type: verifiable.SparseMerkleTreeProof,
	}

	credential, err := isr.newW3CCredential(userID, credentialRaw, expiration, schemaType, credentialsStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new w3c credential")
	}

	coreClaim, err := isr.schemaBuilder.CreateCoreClaim(
		ctx,
		schemaType,
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

	credentialsStatusRaw, err := json.Marshal(credentialsStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal credential status")
	}

	return &data.Claim{
		SchemaURL:        resources.ClaimSchemaList[schemaType].ClaimSchemaURL,
		SchemaType:       schemaType.ToRaw(),
		CoreClaim:        data.NewCoreClaim(coreClaim),
		CredentialStatus: credentialsStatusRaw,
		SignatureProof:   signProof,
		Data:             credentialRaw,
		UserID:           userID.String(),
	}, nil
}

func (isr *issuer) newW3CCredential(
	userID *core.ID,
	credentialRaw []byte,
	expiration *time.Time,
	schemaType resources.ClaimSchemaType,
	credentialStatus verifiable.CredentialStatus,
) (*verifiable.W3CCredential, error) {
	issuanceDate := time.Now()
	credential, err := claims.ParseCredentialSubjectFromSnakeCase(credentialRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse credentials")
	}

	credential.CredentialSubject["id"] = "did:" + userID.String()
	credential.CredentialSubject["type"] = schemaType.ToRaw()
	credential.ID = uuid.NewString()
	credential.Expiration = expiration
	credential.IssuanceDate = &issuanceDate
	credential.Issuer = "did:" + isr.Identifier.String()
	credential.CredentialStatus = credentialStatus
	credential.CredentialSchema = verifiable.CredentialSchema{
		ID:   resources.ClaimSchemaList[schemaType].ClaimSchemaURL,
		Type: verifiable.JSONSchemaValidator2018,
	}
	credential.Context = []string{
		verifiable.JSONLDSchemaW3CCredential2018,
		verifiable.JSONLDSchemaIden3Credential,
		isr.schemaBuilder.CachedSchemas[string(schemaType)].JSONLdContext,
	}
	credential.Type = []string{verifiable.TypeW3CVerifiableCredential, schemaType.ToRaw()}

	if schemaType == resources.QDAOMembershipSchemaType {
		credential.CredentialSubject[claims.IssuanceDateCredentialField] = issuanceDate.Unix()
	}

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
