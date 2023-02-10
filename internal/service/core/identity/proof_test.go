package identity

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"testing"

	"bou.ke/monkey"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state"
	"gitlab.com/q-dev/q-id/issuer/mocks/stub"
)

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
}
