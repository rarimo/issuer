package config

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type StatePublisherConfig struct {
	RetryPeriod time.Duration
}

type statePublisherConfigRaw struct {
	RetryPeriod string `fig:"retry_period,required"` //nolint
}

func (c *config) StatePublisher() *StatePublisherConfig {
	return c.statePublisher.Do(func() interface{} {
		configRaw := statePublisherConfigRaw{}
		err := figure.
			Out(&configRaw).
			From(kv.MustGetStringMap(c.getter, "state_publisher")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		retryPeriod, err := time.ParseDuration(configRaw.RetryPeriod)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse state publisher retry period"))
		}

		return &StatePublisherConfig{
			RetryPeriod: retryPeriod,
		}
	}).(*StatePublisherConfig)
}
