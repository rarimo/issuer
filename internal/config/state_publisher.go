package config

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type StatePublisherConfig struct {
	RetryPeriod   time.Duration `fig:"retry_period,required"`   //nolint
	PublishPeriod time.Duration `fig:"publish_period,required"` //nolint
}

func (c *config) StatePublisher() *StatePublisherConfig {
	return c.statePublisher.Do(func() interface{} {
		cfg := StatePublisherConfig{}
		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "state_publisher")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &cfg
	}).(*StatePublisherConfig)
}
