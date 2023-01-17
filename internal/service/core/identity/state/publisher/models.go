package publisher

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/data"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity/state/publisher/contracts"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/zkp"
)

const (
	pendingQueueLen          = 256
	statePublisherRunnerName = "state_publisher"

	StateDoesntExistErrMsg = "State does not exist"
)

var (
	ErrTransactionFailed               = errors.New("transaction failed")
	ErrOnChainStateIsNotFoundInDB      = errors.New("on chain state is not found in the db")
	ErrInDBMoreThanOneUnprocessedState = errors.New("in db more then 1 unprocessed state after previous session")
)

var (
	invalidSignatureError = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _          = abi.NewType("string", "", nil)
)

type publisher struct {
	log                *logan.Entry
	ethClient          *ethclient.Client
	stateStoreContract *contracts.StateStore
	privateKey         *ecdsa.PrivateKey
	address            common.Address
	chainID            *big.Int
	committedStatesQ   data.CommittedStatesQ

	runnerPeriod time.Duration
	pendingQueue chan *publishedStateInfo
}

type Config struct {
	Log          *logan.Entry
	DB           *pgdb.DB
	EthConfig    *config.EthClientConfig
	RunnerPeriod time.Duration
}

type StateTransitionInfo struct {
	IsOldStateGenesis bool
	Identifier        *core.ID
	LatestState       *merkletree.Hash
	NewState          *merkletree.Hash
	ZKProof           *zkp.ZKProof
}

type contractReadableZKP struct {
	proofA [2]*big.Int
	proofB [2][2]*big.Int
	proofC [2]*big.Int
}

type publishedStateInfo struct {
	CommittedState *data.CommittedState
	Tx             *types.Transaction
}
