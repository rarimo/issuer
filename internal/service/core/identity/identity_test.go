package identity

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"sync"
	"testing"

	"bou.ke/monkey"
	"github.com/ethereum/go-ethereum/common/hexutil"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	"gitlab.com/q-dev/q-id/issuer/mocks/stub"
)

const (
	TestPrivateKey = "0x819b6b1176c547655f9fed5589eaaf1ef4a32aab9b46a4190d13d5c81a822117"
)

func NewTestState(t *testing.T, claimsQStub data.ClaimsQ, committedStateQStub data.CommittedStatesQ) *state.IdentityState {
	// initializing new claims tree (it stores claims issued by the user)
	claimsTree, err := merkletree.NewMerkleTree(context.Background(), memory.NewMemoryStorage(), 64)
	if err != nil {
		assert.Fail(t, "failed to init claims tree: %s", err)
	}

	// initializing new revocation tree (it stores revocation ids of the claims that was revoked)
	revocationsTree, err := merkletree.NewMerkleTree(context.Background(), memory.NewMemoryStorage(), 64)
	if err != nil {
		assert.Fail(t, "failed to init revocations tree: %s", err)
	}

	//initializing new roots tree (it stores the all on-chain published claims-tree roots)
	rootsTree, err := merkletree.NewMerkleTree(context.Background(), memory.NewMemoryStorage(), 64)
	if err != nil {
		assert.Fail(t, "failed to init roots tree: %s", err)
	}

	return &state.IdentityState{
		Mutex:           &sync.Mutex{},
		ClaimsTree:      claimsTree,
		RevocationsTree: revocationsTree,
		RootsTree:       rootsTree,
		ClaimsQ:         claimsQStub,
		CommittedStateQ: committedStateQStub,
	}
}

func ParseBJJPrivateKey(privateKey string) (*babyjub.PrivateKey, error) {
	privateKeyRaw, err := hexutil.Decode(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}

	pk := babyjub.PrivateKey{}
	copy(pk[:], privateKeyRaw)

	return &pk, nil
}

func TestIdentity_generateNewIdentity(t *testing.T) {
	// patching rand.Read to return always the same value
	monkey.Patch(rand.Read, func(b []byte) (int, error) {
		return 0, nil
	})

	type expected struct {
		identifier       string
		authClaimHi      string
		authClaimHv      string
		currentStateHash string
	}

	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		State                *state.IdentityState
	}

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	tests := []struct {
		name       string
		fields     fields
		ctx        context.Context
		wantErr    bool
		beforeTest func(t *testing.T, iden *Identity)
		afterTest  func(t *testing.T, iden *Identity)
	}{
		{
			name:    "everything is ok",
			ctx:     context.Background(),
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(committedState *data.CommittedState) error {
							return nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
			beforeTest: func(t *testing.T, iden *Identity) {
				hi, _ := new(big.Int).SetString("15202137672593991855011800751425397603931319918505833188991062670473684854109", 10)
				hv, _ := new(big.Int).SetString("2351654555892372227640888372176282444150254868378439619268573230312091195718", 10)

				err := iden.State.ClaimsTree.Add(context.Background(), hi, hv)
				assert.Nil(t, err, "failed to add claim to claims tree: %s", err)
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				hi, hv, err := iden.AuthClaim.CoreClaim.HiHv()
				assert.Nil(t, err, "failed to get hi and hv from auth claim: %s", err)

				assert.Equal(t, "11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61", iden.Identifier.String(), "identifier is not correct")
				assert.Equal(t, "15202137672593991855011800751425397603931319918505833188991062670473684854109", hi.String(), "auth claim hi is not correct")
				assert.Equal(t, "2351654555892372227640888372176282444150254868378439619268573230312091195718", hv.String(), "auth claim hv is not correct")
				assert.Equal(t, "69710065...", stateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "incorrect private key",
			ctx:     context.Background(),
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: nil,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, nil
								},
							}
						},
					},
				),
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "53173871...", stateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "failed to insert claim",
			ctx:     context.Background(),
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return errors.New("test error")
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(committedState *data.CommittedState) error {
							return nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to insert genesis state",
			ctx:     context.Background(),
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(committedState *data.CommittedState) error {
							return errors.New("test error")
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, nil
								},
							}
						},
					},
				),
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "69710065...", stateHash.String(), "current state hash is not correct")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				log:                  logan.New().Level(logan.FatalLevel),
				State:                tt.fields.State,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(t, iden)
			}

			if err := iden.generateNewIdentity(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("generateNewIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.afterTest != nil {
				tt.afterTest(t, iden)
			}
		})
	}

	monkey.UnpatchAll()
}

func TestIdentity_parseIdentity(t *testing.T) {
	// patching rand.Read to return always the same value
	monkey.Patch(rand.Read, func(b []byte) (int, error) {
		return 0, nil
	})

	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		State                *state.IdentityState
	}

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	correctAuthCoreClaim, err := claims.NewAuthClaim(correctPrivateKey.Public())
	assert.Nil(t, err, "failed to create auth core claim: %s", err)

	type args struct {
		authClaim       *data.Claim
		genesisStateRaw *data.CommittedState
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		beforeTest func(t *testing.T, iden *Identity)
		afterTest  func(t *testing.T, iden *Identity)
	}{
		{
			name:    "everything is ok",
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
			args: args{
				authClaim: &data.Claim{
					ID: 1,
					CoreClaim: &data.CoreClaim{
						Claim: correctAuthCoreClaim,
					},
				},
				genesisStateRaw: &data.CommittedState{
					ID: 1,
				},
			},
			beforeTest: func(t *testing.T, iden *Identity) {
				hi, _ := new(big.Int).SetString("15202137672593991855011800751425397603931319918505833188991062670473684854109", 10)
				hv, _ := new(big.Int).SetString("2351654555892372227640888372176282444150254868378439619268573230312091195718", 10)

				err := iden.State.ClaimsTree.Add(context.Background(), hi, hv)
				assert.Nil(t, err, "failed to add claim to claims tree: %s", err)
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				hi, hv, err := correctAuthCoreClaim.HiHv()
				assert.Nil(t, err, "failed to get hi and hv from auth claim: %s", err)

				assert.Equal(t, "115zTGHKvFeFLPu3vF9Wx2gBqnxGnzvTpmkHPM2diF", iden.Identifier.String(), "identifier is not correct")
				assert.Equal(t, "15202137672593991855011800751425397603931319918505833188991062670473684854109", hi.String(), "auth claim hi is not correct")
				assert.Equal(t, "2351654555892372227640888372176282444150254868378439619268573230312091195718", hv.String(), "auth claim hv is not correct")
				assert.Equal(t, "69710065...", stateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "core claim is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State:                NewTestState(t, &stub.ClaimsQStub{}, &stub.CommittedStatesQStub{}),
			},
			args: args{
				authClaim: &data.Claim{
					ID:        1,
					CoreClaim: nil,
				},
				genesisStateRaw: &data.CommittedState{
					ID: 1,
				},
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "53173871...", stateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "auth claim is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State:                NewTestState(t, &stub.ClaimsQStub{}, &stub.CommittedStatesQStub{}),
			},
			args: args{
				authClaim: nil,
				genesisStateRaw: &data.CommittedState{
					ID: 1,
				},
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "53173871...", stateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "auth claim is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State:                NewTestState(t, &stub.ClaimsQStub{}, &stub.CommittedStatesQStub{}),
			},
			args: args{
				authClaim: &data.Claim{
					ID: 1,
					CoreClaim: &data.CoreClaim{
						Claim: correctAuthCoreClaim,
					},
				},
				genesisStateRaw: nil,
			},
			afterTest: func(t *testing.T, iden *Identity) {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "53173871...", stateHash.String(), "current state hash is not correct")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				log:                  logan.New().Level(logan.FatalLevel),
				State:                tt.fields.State,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(t, iden)
			}

			if err := iden.parseIdentity(context.Background(), tt.args.authClaim, tt.args.genesisStateRaw); (err != nil) != tt.wantErr {
				t.Errorf("parseIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.afterTest != nil {
				tt.afterTest(t, iden)
			}
		})
	}

	monkey.UnpatchAll()
}

func TestIdentity_saveAuthClaimModel(t *testing.T) {
	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		Identifier           *core.ID
		State                *state.IdentityState
	}

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	correctAuthCoreClaim, err := claims.NewAuthClaim(correctPrivateKey.Public())
	assert.Nil(t, err, "failed to create auth core claim: %s", err)

	correctIdentifier, err := core.IDFromString("11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61")
	assert.Nil(t, err, "failed to create identifier: %s", err)

	type args struct {
		coreAuthClaim *core.Claim
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		beforeTest func(t *testing.T, iden *Identity)
	}{
		{
			name:    "everything is ok",
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
			beforeTest: func(t *testing.T, iden *Identity) {
				hi, _ := new(big.Int).SetString("15202137672593991855011800751425397603931319918505833188991062670473684854109", 10)
				hv, _ := new(big.Int).SetString("2351654555892372227640888372176282444150254868378439619268573230312091195718", 10)

				err := iden.State.ClaimsTree.Add(context.Background(), hi, hv)
				assert.Nil(t, err, "failed to add claim to claims tree: %s", err)
			},
			args: args{
				coreAuthClaim: correctAuthCoreClaim,
			},
		},
		{
			name:    "babu jub jub private key is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: nil,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				coreAuthClaim: correctAuthCoreClaim,
			},
		},
		{
			name:    "core auth claim is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
					},
					&stub.CommittedStatesQStub{},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				coreAuthClaim: nil,
			},
		},
		{
			name:    "core auth claim is nil",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return errors.New("test error")
						},
					},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, nil
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				coreAuthClaim: correctAuthCoreClaim,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				Identifier:           tt.fields.Identifier,
				log:                  logan.New().Level(logan.FatalLevel),
				State:                tt.fields.State,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(t, iden)
			}

			if err := iden.saveAuthClaimModel(context.Background(), tt.args.coreAuthClaim); (err != nil) != tt.wantErr {
				t.Errorf("saveAuthClaimModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	monkey.UnpatchAll()
}

func TestIdentity_GenerateMTP(t *testing.T) {
	// patching rand.Read to return always the same value
	monkey.Patch(rand.Read, func(b []byte) (int, error) {
		return 0, nil
	})

	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		Identifier           *core.ID
		AuthClaim            *data.Claim
		circuitsPath         string
		log                  *logan.Entry
		State                *state.IdentityState
	}

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	correctIdentifier, err := core.IDFromString("11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61")
	assert.Nil(t, err, "failed to parse correct identifier: %s", err)

	coreClaim := &core.Claim{}
	_ = coreClaim.SetIndexDataBytes([]byte("test"), []byte("test"))

	type args struct {
		claim *core.Claim
	}

	type expected struct {
		mtp []byte
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		expected expected
	}{
		{
			name:    "core auth claim is nil",
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										IsGenesis:           true,
										RootsTreeRoot:       nil,
										ClaimsTreeRoot:      []byte{162, 138, 100, 92, 106, 67, 46, 192, 123, 190, 83, 204, 34, 197, 169, 151, 26, 251, 131, 158, 169, 137, 5, 4, 193, 38, 67, 95, 170, 163, 181, 26},
										RevocationsTreeRoot: nil,
									}, nil
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				claim: coreClaim,
			},
			expected: expected{
				mtp: []byte{123, 34, 64, 116, 121, 112, 101, 34, 58, 34, 73, 100, 101, 110, 51, 83, 112, 97, 114, 115, 101, 77, 101, 114, 107, 108, 101, 80, 114, 111, 111, 102, 34, 44, 34, 105, 115, 115, 117, 101, 114, 95, 100, 97, 116, 97, 34, 58, 123, 34, 105, 100, 34, 58, 34, 49, 49, 67, 88, 75, 101, 119, 102, 55, 50, 75, 109, 120, 107, 76, 88, 84, 50, 113, 116, 68, 102, 72, 107, 116, 119, 111, 104, 82, 89, 71, 90, 83, 107, 77, 72, 80, 106, 82, 85, 54, 49, 34, 44, 34, 115, 116, 97, 116, 101, 34, 58, 123, 34, 114, 111, 111, 116, 95, 111, 102, 95, 114, 111, 111, 116, 115, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 99, 108, 97, 105, 109, 115, 95, 116, 114, 101, 101, 95, 114, 111, 111, 116, 34, 58, 34, 97, 50, 56, 97, 54, 52, 53, 99, 54, 97, 52, 51, 50, 101, 99, 48, 55, 98, 98, 101, 53, 51, 99, 99, 50, 50, 99, 53, 97, 57, 57, 55, 49, 97, 102, 98, 56, 51, 57, 101, 97, 57, 56, 57, 48, 53, 48, 52, 99, 49, 50, 54, 52, 51, 53, 102, 97, 97, 97, 51, 98, 53, 49, 97, 34, 44, 34, 114, 101, 118, 111, 99, 97, 116, 105, 111, 110, 95, 116, 114, 101, 101, 95, 114, 111, 111, 116, 34, 58, 34, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 34, 44, 34, 118, 97, 108, 117, 101, 34, 58, 34, 57, 53, 56, 97, 102, 100, 49, 57, 56, 56, 54, 99, 49, 101, 49, 55, 101, 98, 102, 49, 52, 100, 53, 97, 49, 48, 101, 50, 52, 102, 98, 53, 99, 98, 56, 56, 53, 49, 99, 52, 99, 55, 102, 99, 50, 49, 100, 98, 49, 57, 102, 56, 97, 52, 56, 98, 50, 100, 55, 53, 101, 56, 50, 56, 34, 125, 44, 34, 109, 116, 112, 34, 58, 123, 34, 101, 120, 105, 115, 116, 101, 110, 99, 101, 34, 58, 116, 114, 117, 101, 44, 34, 115, 105, 98, 108, 105, 110, 103, 115, 34, 58, 91, 93, 125, 125, 44, 34, 109, 116, 112, 34, 58, 123, 34, 101, 120, 105, 115, 116, 101, 110, 99, 101, 34, 58, 116, 114, 117, 101, 44, 34, 115, 105, 98, 108, 105, 110, 103, 115, 34, 58, 91, 93, 125, 125},
			},
		},
		{
			name:    "claim is absent in the tree",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										IsGenesis:           true,
										RootsTreeRoot:       nil,
										ClaimsTreeRoot:      []byte{16, 138, 100, 92, 106, 67, 46, 192, 123, 190, 83, 204, 34, 197, 169, 151, 26, 251, 131, 158, 169, 137, 5, 4, 193, 38, 67, 95, 170, 163, 181, 26},
										RevocationsTreeRoot: nil,
									}, nil
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				claim: coreClaim,
			},
		},
		{
			name:    "failed to get latest committed state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, errors.New("test error")
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
		},
		{
			name:    "failed to marshal merkle tree proof",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{},
					&stub.CommittedStatesQStub{
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										IsGenesis:           true,
										RootsTreeRoot:       nil,
										ClaimsTreeRoot:      []byte{162, 138, 100, 92, 106, 67, 46, 192, 123, 190, 83, 204, 34, 197, 169, 151, 26, 251, 131, 158, 169, 137, 5, 4, 193, 38, 67, 95, 170, 163, 181, 26},
										RevocationsTreeRoot: nil,
									}, nil
								},
							}
						},
					},
				),
				Identifier: &correctIdentifier,
			},
			args: args{
				claim: coreClaim,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				Identifier:           tt.fields.Identifier,
				circuitsPath:         tt.fields.circuitsPath,
				log:                  tt.fields.log,
				State:                tt.fields.State,
			}

			if tt.name == "failed to marshal merkle tree proof" {
				monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
					return nil, errors.New("test error")
				})
			}

			if tt.args.claim != nil {
				hi, hv, err := tt.args.claim.HiHv()
				assert.Nil(t, err, "failed to get claim hindex: %s", err)

				err = iden.State.ClaimsTree.Add(context.Background(), hi, hv)
				assert.Nil(t, err, "failed to add claim to tree: %s", err)
			}

			actual, err := iden.GenerateMTP(context.Background(), tt.args.claim)
			if (err != nil) != tt.wantErr {
				t.Errorf("iden.GenerateMTP() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expected.mtp != nil {
				assert.Equal(t, tt.expected.mtp, actual, "MTP is not equal")
			}
		})
	}

	monkey.UnpatchAll()
}

func TestIdentity_Init(t *testing.T) {
	// patching rand.Read to return always the same value
	monkey.Patch(rand.Read, func(b []byte) (int, error) {
		return 0, nil
	})

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	authCoreClaim, err := claims.NewAuthClaim(correctPrivateKey.Public())
	assert.Nil(t, err, "failed to create auth core claim: %s", err)

	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		State                *state.IdentityState
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		beforeTest func(t *testing.T, iden *Identity)
		afterTest  func(t *testing.T, iden *Identity)
	}{
		{
			name:    "everything is ok with generating new identity",
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
			afterTest: func(t *testing.T, iden *Identity) {
				currentStateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61", iden.Identifier.String(), "identity is not correct")
				assert.Equal(t, "69710065...", currentStateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "failed to insert auth claim",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return errors.New("test error")
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to get auth claim",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, errors.New("test error")
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to insert committed state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return errors.New("test error")
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to insert committed state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, errors.New("failed to get genesis state")
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to get latest committed state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{}, errors.New("test error")
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "everything is ok without generating new identity",
			wantErr: false,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return &data.Claim{
								CoreClaim: data.NewCoreClaim(authCoreClaim),
							}, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return &data.CommittedState{
								ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
							}, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},

			beforeTest: func(t *testing.T, iden *Identity) {
				hi, _ := new(big.Int).SetString("15202137672593991855011800751425397603931319918505833188991062670473684854109", 10)
				hv, _ := new(big.Int).SetString("2351654555892372227640888372176282444150254868378439619268573230312091195718", 10)

				err := iden.State.ClaimsTree.Add(context.Background(), hi, hv)
				assert.Nil(t, err, "failed to add claim to claims tree: %s", err)
			},
			afterTest: func(t *testing.T, iden *Identity) {
				currentStateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, "11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61", iden.Identifier.String(), "identity is not correct")
				assert.Equal(t, "69710065...", currentStateHash.String(), "current state hash is not correct")
			},
		},
		{
			name:    "failed to insert auth claim",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return errors.New("test error")
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return &data.Claim{
								CoreClaim: data.NewCoreClaim(authCoreClaim),
							}, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return &data.CommittedState{
								ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
							}, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to get auth claim",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return nil, errors.New("test error")
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return &data.CommittedState{
								ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
							}, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to insert committed state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return &data.Claim{
								CoreClaim: data.NewCoreClaim(authCoreClaim),
							}, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return errors.New("test error")
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return &data.CommittedState{
								ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
							}, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to get genesis state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return &data.Claim{
								CoreClaim: data.NewCoreClaim(authCoreClaim),
							}, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return nil, errors.New("test error")
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return &data.CommittedState{
										ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
									}, nil
								},
							}
						},
					},
				),
			},
		},
		{
			name:    "failed to get latest state",
			wantErr: true,
			fields: fields{
				babyJubJubPrivateKey: correctPrivateKey,
				State: NewTestState(t,
					&stub.ClaimsQStub{
						InsertStub: func(claim *data.Claim) error {
							return nil
						},
						GetAuthClaimStub: func() (*data.Claim, error) {
							return &data.Claim{
								CoreClaim: data.NewCoreClaim(authCoreClaim),
							}, nil
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(state *data.CommittedState) error {
							return nil
						},
						GetGenesisStub: func() (*data.CommittedState, error) {
							return &data.CommittedState{
								ClaimsTreeRoot: []byte{235, 103, 37, 219, 83, 72, 67, 31, 122, 226, 128, 79, 67, 109, 249, 63, 122, 24, 185, 51, 54, 244, 182, 211, 240, 209, 86, 183, 5, 238, 134, 8},
							}, nil
						},
						WhereStatusStub: func(status data.Status) data.CommittedStatesQ {
							return &stub.CommittedStatesQStub{
								GetLatestStub: func() (*data.CommittedState, error) {
									return nil, errors.New("test error")
								},
							}
						},
					},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				log:                  logan.New().Level(logan.FatalLevel),
				State:                tt.fields.State,
			}

			if tt.beforeTest != nil {
				tt.beforeTest(t, iden)
			}

			if err := iden.Init(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("iden.Init() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.afterTest != nil {
				tt.afterTest(t, iden)
			}
		})
	}

	monkey.UnpatchAll()
}
