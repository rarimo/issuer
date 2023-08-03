package issuer

import (
	"context"
	"fmt"

	"github.com/iden3/go-jwz"
	"github.com/pkg/errors"

	"gitlab.com/rarimo/identity/issuer/internal/data"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/claims"
	"gitlab.com/rarimo/identity/issuer/internal/service/core/identity"
)

func (isr *issuer) generateProofs(ctx context.Context, claim *data.Claim) (err error) {
	if err != nil {
		return errors.Wrap(err, "failed to create revocation check url")
	}

	issuerData, err := isr.CompactIssuerData(
		ctx,
		fmt.Sprint(isr.baseURL, claims.CredentialStatusCheckURL, isr.AuthClaim.CoreClaim.GetRevocationNonce()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to compact issuer data")
	}

	claim.SignatureProof, err = isr.GenerateIncSigProof(
		claim.CoreClaim.Claim,
		*issuerData,
		fmt.Sprint(isr.baseURL, claims.MTPUpdateURL, isr.AuthClaim.ID),
	)
	if err != nil {
		return errors.Wrap(err, "failed to get signature proof")
	}

	claim.MTP, err = isr.GenerateIncMTProof(
		ctx,
		claim.CoreClaim.Claim,
		fmt.Sprint(isr.baseURL, claims.MTPUpdateURL, claim.ID),
		*issuerData,
	)
	if err != nil {
		if errors.Is(err, identity.ErrClaimWasNotPublishedYet) {
			return nil
		}
		return errors.Wrap(err, "failed to generate merkle tree proof")
	}

	return nil
}

func (isr *issuer) checkClaimRetriever(claim *data.Claim, claimRetriever string, token *jwz.Token) (bool, error) {
	claimRecipientID, err := claim.CoreClaim.GetID()
	if err != nil {
		return false, errors.Wrap(err, "failed to get claim recipient identifier")
	}

	if claimRecipientID.String() != claimRetriever {
		return false, nil
	}

	isZKPValid, err := token.Verify(isr.authVerificationKey)
	if err != nil {
		return false, errors.Wrap(ErrProofVerifyFailed, err.Error())
	}

	return isZKPValid, nil
}

func strptr(str string) *string {
	return &str
}
