package state

import (
	"context"
	"sync"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/pkg/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data/pg"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/publisher"
	treestorage "gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/tree_storage"
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

	statePublisher, err := publisher.NewPublisher(cfg.PublisherConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create publisher")
	}

	go statePublisher.Run(ctx)

	circuits, err := ReadCircuits(cfg.IdentityConfig.CircuitsPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read circuits")
	}

	return &IdentityState{
		circuits:        circuits,
		CommittedStateQ: pg.NewCommittedStateQ(cfg.DB),
		ClaimsQ:         pg.NewClaimsQ(cfg.DB),
		ClaimsTree:      claimsTree,
		RevocationsTree: revsTree,
		RootsTree:       rootsTree,
		publisher:       statePublisher,
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

func (is *IdentityState) SetupGenesis(publicKey *babyjub.PublicKey) (*core.ID, *core.Claim, error) {
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

	identifier, err := core.IdGenesisFromIdenState(core.TypeDefault, genesisStateHash.BigInt())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get identifier from the genesis state")
	}

	return identifier, authClaim, nil
}
