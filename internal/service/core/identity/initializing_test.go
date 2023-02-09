package identity

import (
	"context"
	"crypto/rand"
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
		name     string
		fields   fields
		ctx      context.Context
		wantErr  bool
		expected expected
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
					},
				),
			},
			expected: expected{
				identifier:       "11CXKewf72KmxkLXT2qtDfHktwohRYGZSkMHPjRU61",
				authClaimHi:      "15202137672593991855011800751425397603931319918505833188991062670473684854109",
				authClaimHv:      "2351654555892372227640888372176282444150254868378439619268573230312091195718",
				currentStateHash: "69710065...",
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
						InsertStub: func(committedState *data.CommittedState) error {
							return nil
						},
					},
				),
			},
			expected: expected{
				currentStateHash: "53173871...",
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
							return errors.New("")
						},
					},
					&stub.CommittedStatesQStub{
						InsertStub: func(committedState *data.CommittedState) error {
							return nil
						},
					},
				),
			},
			expected: expected{
				currentStateHash: "69710065...",
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
							return errors.New("")
						},
					},
				),
			},
			expected: expected{
				currentStateHash: "69710065...",
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
			if err := iden.generateNewIdentity(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("generateNewIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expected.identifier != "" {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				hi, hv, err := iden.AuthClaim.CoreClaim.HiHv()
				assert.Nil(t, err, "failed to get hi and hv from auth claim: %s", err)

				assert.Equal(t, tt.expected.identifier, iden.Identifier.String(), "identifier is not correct")
				assert.Equal(t, tt.expected.authClaimHi, hi.String(), "auth claim hi is not correct")
				assert.Equal(t, tt.expected.authClaimHv, hv.String(), "auth claim hv is not correct")
				assert.Equal(t, tt.expected.currentStateHash, stateHash.String(), "current state hash is not correct")
			}

			if tt.expected.currentStateHash != "" {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, tt.expected.currentStateHash, stateHash.String(), "current state hash is not correct")
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

	type expected struct {
		identifier       string
		authClaimHi      string
		authClaimHv      string
		currentStateHash string
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
		name     string
		fields   fields
		args     args
		wantErr  bool
		expected expected
	}{
		{
			name:    "everything is ok",
			wantErr: false,
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
				genesisStateRaw: &data.CommittedState{
					ID: 1,
				},
			},
			expected: expected{
				identifier:       "115zTGHKvFeFLPu3vF9Wx2gBqnxGnzvTpmkHPM2diF",
				authClaimHi:      "15202137672593991855011800751425397603931319918505833188991062670473684854109",
				authClaimHv:      "2351654555892372227640888372176282444150254868378439619268573230312091195718",
				currentStateHash: "53173871...",
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
			expected: expected{
				currentStateHash: "53173871...",
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
			expected: expected{
				currentStateHash: "53173871...",
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
			expected: expected{
				currentStateHash: "53173871...",
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
			if err := iden.parseIdentity(context.Background(), tt.args.authClaim, tt.args.genesisStateRaw); (err != nil) != tt.wantErr {
				t.Errorf("parseIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.expected.identifier != "" {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				hi, hv, err := iden.AuthClaim.CoreClaim.HiHv()
				assert.Nil(t, err, "failed to get hi and hv from auth claim: %s", err)

				assert.Equal(t, tt.expected.identifier, iden.Identifier.String(), "identifier is not correct")
				assert.Equal(t, tt.expected.authClaimHi, hi.String(), "auth claim hi is not correct")
				assert.Equal(t, tt.expected.authClaimHv, hv.String(), "auth claim hv is not correct")
				assert.Equal(t, tt.expected.currentStateHash, stateHash.String(), "current state hash is not correct")
			}

			if tt.expected.currentStateHash != "" {
				stateHash, err := iden.State.GetCurrentStateHash()
				assert.Nil(t, err, "failed to get current state hash: %s", err)

				assert.Equal(t, tt.expected.currentStateHash, stateHash.String(), "current state hash is not correct")
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
					&stub.CommittedStatesQStub{},
				),
				Identifier: &correctIdentifier,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iden := &Identity{
				babyJubJubPrivateKey: tt.fields.babyJubJubPrivateKey,
				Identifier:           tt.fields.Identifier,
				log:                  logan.New().Level(logan.FatalLevel),
				State:                tt.fields.State,
			}
			if err := iden.saveAuthClaimModel(context.Background(), tt.args.coreAuthClaim); (err != nil) != tt.wantErr {
				t.Errorf("saveAuthClaimModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
