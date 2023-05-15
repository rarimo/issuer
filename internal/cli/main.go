package cli

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/q-dev/q-id/issuer/internal/config"
	"gitlab.com/q-dev/q-id/issuer/internal/service/api"
	"gitlab.com/q-dev/q-id/issuer/internal/service/core/identity"
)

func Run(args []string) bool {
	log := logan.New()
	setRandSeed()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	cfg := config.New(kv.MustFromEnv())
	log = cfg.Log()

	app := kingpin.New("issuer", "")

	runCmd := app.Command("run", "run command")
	serviceCmd := runCmd.Command("service", "run service")

	migrateCmd := app.Command("migrate", "migrate command")
	migrateUpCmd := migrateCmd.Command("up", "migrate db up")
	migrateDownCmd := migrateCmd.Command("down", "migrate db down")

	initCmd := app.Command("init", "init command")
	identityCmd := initCmd.Command("identity", "initialize new issuer's identity state if not present")

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch cmd {
	case serviceCmd.FullCommand():
		run(wg, ctx, cfg, api.Run)
	case migrateUpCmd.FullCommand():
		err = MigrateUp(cfg)
	case migrateDownCmd.FullCommand():
		err = MigrateDown(cfg)
	case identityCmd.FullCommand():
		_, err = identity.New(ctx, cfg)
	default:
		log.Errorf("unknown command %s", cmd)
		cancel()
		return false
	}
	if err != nil {
		log.WithError(err).Error("failed to exec cmd")
		cancel()
		return false
	}

	graceful := make(chan os.Signal, 1)
	signal.Notify(graceful, os.Interrupt, syscall.SIGTERM)

	waitGroupChannel := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitGroupChannel)
	}()

	select {
	case <-waitGroupChannel:
		log.Info("all services stopped")
		cancel()
	case <-graceful:
		log.Info("got signal to stop gracefully")
		cancel()
		<-waitGroupChannel
	}

	return true
}

func run(
	wg *sync.WaitGroup, ctx context.Context,
	cfg config.Config, runner func(context.Context, config.Config),
) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		runner(ctx, cfg)
	}()
}

func setRandSeed() {
	var seed [8]byte
	_, err := cryptorand.Read(seed[:])
	if err != nil {
		panic(err)
	}
	mathrand.Seed(int64(binary.LittleEndian.Uint64(seed[:])))
}
