package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	types.Copuser

	EthClient() *EthClientConfig
	StatePublisher() *StatePublisherConfig
	Identity() *IdentityConfig
	Issuer() *IssuerConfig
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	types.Copuser
	getter kv.Getter

	ethClient      comfig.Once
	statePublisher comfig.Once
	issuer         comfig.Once
	identity       comfig.Once
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
