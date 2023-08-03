package issuer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"

	"github.com/rarimo/issuer/internal/data"
	"github.com/rarimo/issuer/internal/service/core/claims"
)

func (isr *issuer) compactClaim(
	ctx context.Context,
	userDID *core.DID,
	expiration *time.Time,
	claimType claims.ClaimSchemaType,
	credentialSubjectRaw []byte,
) (*data.Claim, error) {
	revNonce, err := claims.CryptoRandUint64()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random uint64")
	}

	credentialsStatus := verifiable.CredentialStatus{
		ID:              fmt.Sprint(isr.baseURL, claims.CredentialStatusCheckURL, revNonce),
		Type:            verifiable.SparseMerkleTreeProof,
		RevocationNonce: revNonce,
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

	return &data.Claim{
		ID:         claimID,
		ClaimType:  claimType.ToRaw(),
		CoreClaim:  data.NewCoreClaim(coreClaim),
		Credential: credentialRaw,
		UserID:     userDID.ID.String(),
	}, nil
}

func (isr *issuer) newW3CCredential(
	claimID string,
	userDID *core.DID,
	credentialSubjectRaw []byte,
	expiration *time.Time,
	claimType claims.ClaimSchemaType,
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
		ID:   fmt.Sprint(isr.schemaBuilder.SchemasBaseURL, claims.ClaimSchemaList[claimType].ClaimSchemaURL),
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

func NewClaimOffer(callBackURL string, from, to *core.DID, claim *data.Claim) *protocol.CredentialsOfferMessage {
	return &protocol.CredentialsOfferMessage{
		ID:       uuid.NewString(),
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.CredentialOfferMessageType,
		ThreadID: uuid.NewString(),
		Body: protocol.CredentialsOfferMessageBody{
			URL: callBackURL,
			Credentials: []protocol.CredentialOffer{
				protocol.CredentialOffer{
					ID:          claim.ID,
					Description: claims.ClaimSchemaList[claims.ClaimSchemaTypeList[claim.ClaimType]].ClaimSchemaName,
				},
			},
		},
		From: from.String(),
		To:   to.String(),
	}
}

func ClaimOfferToRaw(
	claimOffer *protocol.CredentialsOfferMessage,
	createdAt time.Time,
	fromID core.ID,
	toID core.ID,
) *data.ClaimOffer {
	claimOfferRaw := data.ClaimOffer{
		ID:         claimOffer.ThreadID,
		From:       fromID.String(),
		To:         toID.String(),
		CreatedAt:  createdAt,
		IsReceived: false,
	}

	if len(claimOffer.Body.Credentials) > 0 {
		claimOfferRaw.ClaimID = claimOffer.Body.Credentials[0].ID
	}

	return &claimOfferRaw
}

func ClaimModelToW3Credential(claim *data.Claim) (*verifiable.W3CCredential, error) {
	res := &verifiable.W3CCredential{
		Proof: verifiable.CredentialProofs{
			claim.SignatureProof,
			claim.MTP,
		},
	}
	err := json.Unmarshal(claim.Credential, res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credential")
	}

	return res, nil
}
