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
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-rapidsnark/prover"
	"github.com/iden3/go-rapidsnark/witness"
	"github.com/pkg/errors"

	dataPkg "gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/zkp"
)

func (is *IdentityState) GenerateStateCommitment(
	ctx context.Context,
) (*StateTransitionInfo, *dataPkg.CommittedState, error) {
	is.Lock()
	defer is.Unlock()

	processingStates, err := is.DB.CommittedStatesQ().WhereStatus(dataPkg.StatusProcessing).Select()
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to select committed states")
	}
	if len(processingStates) > 0 {
		return nil, nil, nil
	}

	transitionInfo, err := is.PrepareTransitionInfo(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to prepare state transition info")
	}

	newStateCommit := (&CommittedState{
		Status:              dataPkg.StatusProcessing,
		CreatedAt:           time.Now(),
		IsGenesis:           false,
		RootsTreeRoot:       is.RootsTree.Root(),
		ClaimsTreeRoot:      is.ClaimsTree.Root(),
		RevocationsTreeRoot: is.RevocationsTree.Root(),
	}).ToRaw()

	err = is.DB.CommittedStatesQ().Insert(newStateCommit)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to insert committed state into db")
	}

	return transitionInfo, newStateCommit, nil
}

func (is *IdentityState) PrepareTransitionInfo(ctx context.Context) (*StateTransitionInfo, error) {
	oldStateRaw, err := is.DB.CommittedStatesQ().New().WhereStatus(dataPkg.StatusCompleted).GetLatest()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the last committed state from db")
	}
	if oldStateRaw == nil {
		return nil, ErrOldStateNotFound
	}

	oldState, err := CommittedStateFromRaw(oldStateRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the latest committed state")
	}

	oldStateHash, err := oldState.StateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest committed state hash")
	}

	err = is.RootsTree.Add(ctx, oldState.ClaimsTreeRoot.BigInt(), merkletree.HashZero.BigInt())
	if err != nil && !errors.Is(err, merkletree.ErrEntryIndexAlreadyExists) {
		return nil, errors.Wrap(err, "failed to add new claim tree root to the roots tree")
	}

	newStateHash, err := is.GetCurrentStateHash()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the new state hash")
	}

	if oldStateHash.Equals(newStateHash) ||
		(oldState.ClaimsTreeRoot.Equals(is.ClaimsTree.Root()) &&
			oldState.RevocationsTreeRoot.Equals(is.RevocationsTree.Root())) {
		return nil, ErrStateWasntChanged
	}

	transitionInputs, err := is.prepareTransitionInputs(ctx, oldState, newStateHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare state transition inputs")
	}

	transitionProof, err := is.generateTransitionProof(transitionInputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate state transition proof")
	}

	return &StateTransitionInfo{
		Identifier:        is.Identifier,
		LatestState:       oldStateHash,
		NewState:          newStateHash,
		IsOldStateGenesis: oldState.IsGenesis,
		ZKProof:           transitionProof.Proof,
	}, nil
}

func (is *IdentityState) prepareTransitionInputs(
	ctx context.Context,
	oldState *CommittedState,
	newStateHash *merkletree.Hash,
) ([]byte, error) {
	oldStateCircuits, err := circuitsState(oldState)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get circuit state")
	}

	newStateCircuits := circuits.TreeState{
		State:          newStateHash,
		ClaimsRoot:     is.ClaimsTree.Root(),
		RevocationRoot: is.RevocationsTree.Root(),
		RootOfRoots:    is.RootsTree.Root(),
	}

	authInclusionProof, err := is.GetInclusionProof(ctx, is.AuthClaim, oldState.ClaimsTreeRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the auth claim inclusion proof")
	}

	authNewInclusionProof, err := is.GetInclusionProof(ctx, is.AuthClaim, newStateCircuits.ClaimsRoot)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the auth claim new inclusion proof")
	}

	authNonRevocationProof, _, err := is.RevocationsTree.GenerateProof(
		ctx,
		big.NewInt(int64(is.AuthClaim.GetRevocationNonce())),
		oldState.RevocationsTreeRoot,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the auth claim non revocation proof")
	}

	hashOldAndNewStates, err := poseidon.Hash([]*big.Int{oldStateCircuits.State.BigInt(), newStateHash.BigInt()})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get poseidon hash")
	}

	transitionInputs, err := circuits.StateTransitionInputs{
		ID:                      is.Identifier,
		NewTreeState:            newStateCircuits,
		OldTreeState:            oldStateCircuits,
		IsOldStateGenesis:       oldState.IsGenesis,
		AuthClaim:               is.AuthClaim,
		AuthClaimIncMtp:         authInclusionProof,
		AuthClaimNonRevMtp:      authNonRevocationProof,
		AuthClaimNewStateIncMtp: authNewInclusionProof,
		Signature:               is.BabyJubJubPrivateKey.SignPoseidon(hashOldAndNewStates),
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

func circuitsState(committedState *CommittedState) (circuits.TreeState, error) {
	stateHash, err := committedState.StateHash()
	if err != nil {
		return circuits.TreeState{}, errors.Wrap(err, "failed to get state hash")
	}

	return circuits.TreeState{
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
