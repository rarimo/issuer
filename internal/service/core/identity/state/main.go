package state

import (
	"context"
	"sync"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/pkg/errors"

	"github.com/rarimo/issuer/internal/data/pg"
	"github.com/rarimo/issuer/internal/service/core/claims"
	treestorage "github.com/rarimo/issuer/internal/service/core/identity/state/tree_storage"
)

func NewIdentityState(ctx context.Context, cfg Config) (*IdentityState, error) {
	claimsTree, err := merkletree.NewMerkleTree(
		ctx, treestorage.NewTreeStorage(cfg.DB, treestorage.ClaimsTreeName), cfg.IdentityConfig.TreeDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new claims merkle tree")
	}

	revsTree, err := merkletree.NewMerkleTree(
		ctx, treestorage.NewTreeStorage(cfg.DB, treestorage.RevocationTreeName), cfg.IdentityConfig.TreeDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new revocations merkle tree")
	}

	rootsTree, err := merkletree.NewMerkleTree(
		ctx, treestorage.NewTreeStorage(cfg.DB, treestorage.RootsTreeName), cfg.IdentityConfig.TreeDepth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new roots merkle tree")
	}

	circuits, err := ReadCircuits(cfg.IdentityConfig.CircuitsPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read circuits")
	}

	return &IdentityState{
		circuits:        circuits,
		DB:              pg.NewMasterQ(cfg.DB),
		ClaimsTree:      claimsTree,
		RevocationsTree: revsTree,
		RootsTree:       rootsTree,
		Mutex:           &sync.Mutex{},
	}, nil
}

func ReadCircuits(path string) (circuits map[string][]byte, err error) {
	circuits = make(map[string][]byte)

	circuits[StateTransitionCircuitFinalKey], err = ReadFileByPath(path, StateTransitionCircuitFinalKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read state transition circuit final")
	}

	circuits[StateTransitionCircuitWasmPath], err = ReadFileByPath(path, StateTransitionCircuitWasmPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read state transition circuit wasm")
	}

	return circuits, nil
}

func (is *IdentityState) GetCurrentStateHash() (*merkletree.Hash, error) {
	hash, err := merkletree.HashElems(
		is.ClaimsTree.Root().BigInt(),
		is.RevocationsTree.Root().BigInt(),
		is.RootsTree.Root().BigInt(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hash merkle trees roots")
	}

	return hash, nil
}

func (is *IdentityState) AddClaimMT(claim *core.Claim) error {
	is.Lock()
	defer is.Unlock()

	indexHash, valueHash, err := claim.HiHv()
	if err != nil {
		return errors.Wrap(err, "failed to get hash of the index and value from the claim")
	}

	err = is.ClaimsTree.Add(context.Background(), indexHash, valueHash)
	if err != nil {
		return errors.Wrap(err, "failed to add leaf to merkle tree")
	}

	return nil
}

func (is *IdentityState) SetupGenesis(publicKey *babyjub.PublicKey) (*core.DID, *core.Claim, error) {
	authClaim, err := claims.NewAuthClaim(publicKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to crate new auth claim")
	}

	err = is.AddClaimMT(authClaim)
	if err != nil && !errors.Is(err, merkletree.ErrEntryIndexAlreadyExists) {
		return nil, nil, errors.Wrap(err, "failed to add the claim to merkle tree")
	}

	genesisStateHash, err := is.GetCurrentStateHash()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get state hash")
	}

	didType, err := core.BuildDIDType(core.DIDMethodIden3, core.NoChain, core.NoNetwork)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build did type")
	}

	did, err := core.DIDGenesisFromIdenState(didType, genesisStateHash.BigInt())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get did from the genesis state")
	}

	return did, authClaim, nil
}

func (is *IdentityState) SetIdentityInfo(identityInfo *IdentityInfo) {
	is.IdentityInfo = identityInfo
}
