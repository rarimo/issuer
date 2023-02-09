package state

import (
	"context"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"
	dataPkg "gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/publisher"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/zkp"
)

func (is *IdentityState) PublishOnChain(ctx context.Context, identityInfo *IdentityInfo) (string, error) {
	is.Lock()
	defer is.Unlock()

	processingStates, err := is.CommittedStateQ.New().WhereStatus(dataPkg.StatusProcessing).Select()
	if err != nil {
		return "", errors.Wrap(err, "failed to select committed states")
	}
	if len(processingStates) > 0 {
		return processingStates[0].TxID, nil
	}

	transitionInfo, err := is.prepareTransitionInfo(ctx, identityInfo)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare state transition info")
	}

	newStateCommit := (&CommittedState{
		Status:              dataPkg.StatusProcessing,
		CreatedAt:           time.Now(),
		IsGenesis:           false,
		RootsTreeRoot:       is.RootsTree.Root(),
		ClaimsTreeRoot:      is.ClaimsTree.Root(),
		RevocationsTreeRoot: is.RevocationsTree.Root(),
	}).ToRaw()

	err = is.CommittedStateQ.Insert(newStateCommit)
	if err != nil {
		return "", errors.Wrap(err, "failed to insert committed state into db")
	}

	txHash, err := is.publisher.PublishState(ctx, transitionInfo, newStateCommit)
	if err != nil {
		return "", errors.Wrap(err, "failed to publish state with publisher")
	}

	return txHash, nil
}

func (is *IdentityState) prepareTransitionInfo(
	ctx context.Context,
	identityInfo *IdentityInfo,
) (*publisher.StateTransitionInfo, error) {
	latestStateRaw, err := is.CommittedStateQ.New().WhereStatus(dataPkg.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the last committed state from db")
	}
	if latestStateRaw == nil {
		return nil, errors.New("the latest state is absent")
	}

	latestState, err := CommittedStateFromRaw(latestStateRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the latest committed state")
	}

	transitionInputs, err := is.prepareTransitionInputs(ctx, identityInfo, latestState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare state transition inputs")
	}

	latestStateHash, err := latestState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state hash")
	}

	newStateHash, err := is.GetCurrentStateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the new state hash")
	}

	if latestStateHash.Equals(newStateHash) ||
		(latestState.ClaimsTreeRoot.Equals(is.ClaimsTree.Root()) &&
			latestState.RevocationsTreeRoot.Equals(is.RevocationsTree.Root())) {
		return nil, ErrStateWasntChanged
	}

	transitionProof, err := is.generateTransitionProof(transitionInputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate state transition proof")
	}

	return &publisher.StateTransitionInfo{
		Identifier:        identityInfo.Identifier,
		LatestState:       latestStateHash,
		NewState:          newStateHash,
		IsOldStateGenesis: latestState.IsGenesis,
		ZKProof:           transitionProof.Proof,
	}, nil
}

func (is *IdentityState) prepareTransitionInputs(
	ctx context.Context,
	identityInfo *IdentityInfo,
	latestState *CommittedState,
) ([]byte, error) {
	oldState, err := circuitsState(latestState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get circuit state")
	}

	err = is.RootsTree.Add(ctx, latestState.ClaimsTreeRoot.BigInt(), merkletree.HashZero.BigInt())
	if err != nil && !errors.Is(err, merkletree.ErrEntryIndexAlreadyExists) {
		return nil, errors.Wrap(err, "failed to add new claim tree root to the roots tree")
	}

	newStateHash, err := is.GetCurrentStateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current state hash")
	}

	authClaimInclusionProof, err := is.GetInclusionProof(ctx, identityInfo.AuthClaim, latestState.ClaimsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the auth claim inclusion proof")
	}

	authNonRevocationProof, err := is.GetRevocationProof(ctx, identityInfo.AuthClaim, latestState.RevocationsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the auth claim non revocation proof")
	}

	authClaimWithProofs := circuits.Claim{
		Claim:     identityInfo.AuthClaim,
		TreeState: *oldState,
		Proof:     authClaimInclusionProof,
		NonRevProof: &circuits.ClaimNonRevStatus{
			TreeState: *oldState,
			Proof:     authNonRevocationProof,
		},
	}

	hashOldAndNewStates, err := poseidon.Hash([]*big.Int{oldState.State.BigInt(), newStateHash.BigInt()})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get poseidon hash")
	}

	transitionInputs, err := circuits.StateTransitionInputs{
		ID:                identityInfo.Identifier,
		NewState:          newStateHash,
		OldTreeState:      *oldState,
		IsOldStateGenesis: latestState.IsGenesis,
		AuthClaim:         authClaimWithProofs,
		Signature:         identityInfo.BabyJubJubPrivateKey.SignPoseidon(hashOldAndNewStates),
	}.InputsMarshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal state transition inputs")
	}

	return transitionInputs, nil
}

func (is *IdentityState) GetInclusionProof(
	ctx context.Context,
	claim *core.Claim,
	claimTreeRoot *merkletree.Hash,
) (*merkletree.Proof, error) {
	hi, err := claim.HIndex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hash of index and a value of the claim")
	}

	proof, _, err := is.ClaimsTree.GenerateProof(ctx, hi, claimTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a proof")
	}

	return proof, nil
}

func (is *IdentityState) GetRevocationProof(
	ctx context.Context,
	claim *core.Claim,
	revocationsTreeRoot *merkletree.Hash,
) (*merkletree.Proof, error) {
	hi, err := claim.HIndex()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hash of index of the claim")
	}

	proof, _, err := is.RevocationsTree.GenerateProof(ctx, hi, revocationsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a proof")
	}

	return proof, nil
}

func circuitsState(committedState *CommittedState) (*circuits.TreeState, error) {
	stateHash, err := committedState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get state hash")
	}

	return &circuits.TreeState{
		State:          stateHash,
		ClaimsRoot:     committedState.ClaimsTreeRoot,
		RevocationRoot: committedState.RevocationsTreeRoot,
		RootOfRoots:    committedState.RootsTreeRoot,
	}, nil
}

func (is *IdentityState) generateTransitionProof(transitionInputs []byte) (*zkp.FullProof, error) {
	wtnsCalculator, err := witness.NewCircom2WitnessCalculator(is.circuits[StateTransitionCircuitWasmPath], true)
	if err != nil {
		return nil, errors.New("failed to create witness calculator")
	}

	parsedInputs, err := witness.ParseInputs(transitionInputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse witness inputs")
	}

	wtnsBytes, err := wtnsCalculator.CalculateWTNSBin(parsedInputs, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate witnesses")
	}

	rapidProof, err := prover.Groth16Prover(is.circuits[StateTransitionCircuitFinalKey], wtnsBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate prov with groth16")
	}

	return &zkp.FullProof{
		Proof: &zkp.ZKProof{
			A:        rapidProof.Proof.A,
			B:        rapidProof.Proof.B,
			C:        rapidProof.Proof.C,
			Protocol: rapidProof.Proof.Protocol,
		},
		PubSignals: rapidProof.PubSignals,
	}, nil
}

func (cs *CommittedState) StateHash() (*merkletree.Hash, error) {
	hash, err := merkletree.HashElems(
		cs.ClaimsTreeRoot.BigInt(),
		cs.RevocationsTreeRoot.BigInt(),
		cs.RootsTreeRoot.BigInt(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hash merkle trees roots")
	}

	return hash, nil
}

func ReadFileByPath(basePath string, fileName string) ([]byte, error) {
	path := filepath.Join(basePath, fileName)
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file '%s' by path '%s'", fileName, path)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file '%s' by path '%s'", fileName, path)
	}
	return data, nil
}
