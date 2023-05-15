package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IssuerConfig struct {
	BaseURL        string `fig:"base_url,required"`
	SchemasBaseURL string `fig:"schemas_base_url,required"`
}

func (c *config) Issuer() *IssuerConfig {
	return c.issuer.Do(func() interface{} {
		cfg := IssuerConfig{}
		err := figure.
			Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "issuer")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return &cfg
	}).(*IssuerConfig)
}
