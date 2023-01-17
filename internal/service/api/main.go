package api

import (
	"context"
	"net"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/q-dev/q-id/issuer/internal/config"
	issuerPkg "gitlab.com/q-dev/q-id/issuer/internal/service/core/issuer"
)

type service struct {
	log      *logan.Entry
	listener net.Listener
	copus    types.Copus
	issuer   issuerPkg.Issuer
}

func newService(ctx context.Context, cfg config.Config) (*service, error) {
	issuer, err := issuerPkg.New(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new issuer")
	}

	return &service{
		log:      cfg.Log().WithField("service", "api"),
		listener: cfg.Listener(),
		copus:    cfg.Copus(),
		issuer:   issuer,
	}, nil
}

func Run(ctx context.Context, cfg config.Config) {
	svc, err := newService(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if err := svc.copus.RegisterChi(svc.router()); err != nil {
		panic(errors.Wrap(err, "cop failed"))
	}

	svc.log.Info("Service started")
	ape.Serve(ctx, svc.router(), cfg, ape.ServeOpts{})
}
