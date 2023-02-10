package identity

import (
	"context"
	"crypto/rand"
	"math/big"
	"sync"
	"testing"

	"bou.ke/monkey"
	"github.com/ethereum/go-ethereum/common/hexutil"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	"gitlab.com/q-dev/q-id/issuer/mocks/stub"
)

const (
	TestPrivateKey          = "0x819b6b1176c547655f9fed5589eaaf1ef4a32aab9b46a4190d13d5c81a822117"
	TestIncorrectPrivateKey = "xyz"
)

func NewTestIdentity(t *testing.T) *Identity {
	privateKeyRaw, err := hexutil.Decode(TestPrivateKey)
	if err != nil {
		assert.Fail(t, "failed to decode private key: %s", err)
	}

	var privateKey babyjub.PrivateKey
	copy(privateKey[:], privateKeyRaw)

	return &Identity{
		babyJubJubPrivateKey: &privateKey,
		circuitsPath:         "./circuits",
		State:                NewTestState(t, nil, nil),
	}
}

func NewAuthCoreClaim(publicKey *babyjub.PublicKey) (*core.Claim, error) {
	schemaHash, err := core.NewSchemaHashFromHex(claims.AuthBJJCredentialHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load the schema hash from hex")
	}

	authClaim, err := claims.NewAuthClaim(publicKey, schemaHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed to crate new auth claim")
	}

	return authClaim, nil
}

func NewTestState(t *testing.T, claimsQStub data.ClaimsQ, committedStateQStub data.CommittedStatesQ) *state.IdentityState {
	memStorage := memory.NewMemoryStorage()

	// initializing new claims tree (it stores claims issued by the user)
	claimsTree, err := merkletree.NewMerkleTree(context.Background(), memStorage.WithPrefix([]byte("claims")), 64)
	if err != nil {
		assert.Fail(t, "failed to init claims tree: %s", err)
	}

	// initializing new revocation tree (it stores revocation ids of the claims that was revoked)
	revocationsTree, err := merkletree.NewMerkleTree(context.Background(), memStorage.WithPrefix([]byte("revocations")), 64)
	if err != nil {
		assert.Fail(t, "failed to init revocations tree: %s", err)
	}

	//initializing new roots tree (it stores the all on-chain published claims-tree roots)
	rootsTree, err := merkletree.NewMerkleTree(context.Background(), memStorage.WithPrefix([]byte("roots")), 64)
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

	var incorrectPrivateKey *babyjub.PrivateKey
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
				babyJubJubPrivateKey: incorrectPrivateKey,
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

	correctAuthCoreClaim, err := NewAuthCoreClaim(correctPrivateKey.Public())
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
}

func TestIdentity_saveAuthClaimModel(t *testing.T) {
	type fields struct {
		babyJubJubPrivateKey *babyjub.PrivateKey
		Identifier           *core.ID
		State                *state.IdentityState
	}

	correctPrivateKey, err := ParseBJJPrivateKey(TestPrivateKey)
	assert.Nil(t, err, "failed to parse correct private key: %s", err)

	correctAuthCoreClaim, err := NewAuthCoreClaim(correctPrivateKey.Public())
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
}
