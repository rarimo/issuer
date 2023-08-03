package claims

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	mt "github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-schema-processor/verifiable"
)

// This is forked from github.com/iden3/go-schema-processor/verifiable/proofs.go
// Additional fields was added to the BJJSignatureProof2021 and Iden3SparseMerkleTreeProof structs

// BJJSignatureProof2021 JSON-LD BBJJSignatureProof
type BJJSignatureProof2021 struct {
	Type                 verifiable.ProofType  `json:"type"`
	IssuerProofUpdateURL string                `json:"issuerProofUpdateUrl"`
	IssuerData           verifiable.IssuerData `json:"issuerData"`
	CoreClaim            string                `json:"coreClaim"`
	Signature            string                `json:"signature"`
}

func (p *BJJSignatureProof2021) UnmarshalJSON(in []byte) error {
	var obj struct {
		Type       verifiable.ProofType `json:"type"`
		IssuerData json.RawMessage      `json:"issuerData"`
		CoreClaim  string               `json:"coreClaim"`
		Signature  string               `json:"signature"`
	}
	err := json.Unmarshal(in, &obj)
	if err != nil {
		return err
	}
	if obj.Type != verifiable.BJJSignatureProofType {
		return errors.New("invalid proof type")
	}
	p.Type = obj.Type
	err = json.Unmarshal(obj.IssuerData, &p.IssuerData)
	if err != nil {
		return err
	}
	if err := validateHexCoreClaim(obj.CoreClaim); err != nil {
		return err
	}
	p.CoreClaim = obj.CoreClaim
	if err := validateCompSignature(obj.Signature); err != nil {
		return err
	}
	p.Signature = obj.Signature
	return nil
}

func validateHexCoreClaim(in string) error {
	var claim core.Claim
	err := claim.FromHex(in)
	return err
}

func validateCompSignature(in string) error {
	sigBytes, err := hex.DecodeString(in)
	if err != nil {
		return err
	}
	var sig babyjub.SignatureComp
	if len(sigBytes) != len(sig) {
		return errors.New("invalid signature length")
	}
	copy(sig[:], sigBytes)
	_, err = sig.Decompress()
	return err
}

func (p *BJJSignatureProof2021) ProofType() verifiable.ProofType {
	return p.Type
}

func (p *BJJSignatureProof2021) GetCoreClaim() (*core.Claim, error) {
	var coreClaim core.Claim
	err := coreClaim.FromHex(p.CoreClaim)
	return &coreClaim, err
}

// Iden3SparseMerkleTreeProof JSON-LD structure
type Iden3SparseMerkleTreeProof struct {
	ID   string               `json:"id"`
	Type verifiable.ProofType `json:"type"`

	IssuerData verifiable.IssuerData `json:"issuerData"`
	CoreClaim  string                `json:"coreClaim"`

	MTP *mt.Proof `json:"mtp"`
}

func (p *Iden3SparseMerkleTreeProof) UnmarshalJSON(in []byte) error {
	var obj struct {
		Type       verifiable.ProofType `json:"type"`
		IssuerData json.RawMessage      `json:"issuerData"`
		CoreClaim  string               `json:"coreClaim"`
		MTP        *mt.Proof            `json:"mtp"`
	}
	err := json.Unmarshal(in, &obj)
	if err != nil {
		return err
	}
	if obj.Type != verifiable.Iden3SparseMerkleTreeProofType {
		return errors.New("invalid proof type")
	}
	p.Type = obj.Type
	err = json.Unmarshal(obj.IssuerData, &p.IssuerData)
	if err != nil {
		return err
	}
	if err := validateHexCoreClaim(obj.CoreClaim); err != nil {
		return err
	}
	p.CoreClaim = obj.CoreClaim
	p.MTP = obj.MTP
	return nil
}

func (p *Iden3SparseMerkleTreeProof) ProofType() verifiable.ProofType {
	return p.Type
}

func (p *Iden3SparseMerkleTreeProof) GetCoreClaim() (*core.Claim, error) {
	var coreClaim core.Claim
	err := coreClaim.FromHex(p.CoreClaim)
	return &coreClaim, err
}
