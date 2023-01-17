package config

import (
	"encoding/hex"

	"github.com/iden3/go-iden3-crypto/babyjub"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type IdentityConfig struct {
	CircuitsPath         string
	TreeDepth            int
	BabyJubJubPrivateKey *babyjub.PrivateKey
}

type identityConfig struct {
	CircuitsPath         string `fig:"circuits_path,required"`
	TreeDepth            int    `fig:"tree_depth,required"`
	BabyJubJubPrivateKey string `fig:"babyjubjub_private_key,required"`
}

func (c *config) Identity() *IdentityConfig {
	return c.identity.Do(func() interface{} {
		cfgRaw := identityConfig{}
		err := figure.
			Out(&cfgRaw).
			From(kv.MustGetStringMap(c.getter, "identity")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		parsedConfig, err := parseConfig(&cfgRaw)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse config"))
		}

		return parsedConfig
	}).(*IdentityConfig)
}

func parseConfig(configRaw *identityConfig) (*IdentityConfig, error) {
	babyJubJubPrivateKey, err := parseBabyJubJubPrivateKey(configRaw.BabyJubJubPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse BabyJubJub private key")
	}

	return &IdentityConfig{
		CircuitsPath:         configRaw.CircuitsPath,
		TreeDepth:            configRaw.TreeDepth,
		BabyJubJubPrivateKey: babyJubJubPrivateKey,
	}, nil
}

func parseBabyJubJubPrivateKey(privateKeyRaw string) (*babyjub.PrivateKey, error) {
	if len(privateKeyRaw) == 0 {
		return nil, errors.New("BabyJubJub private key is absent")
	}

	privateKey, err := hex.DecodeString(privateKeyRaw)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode from hex")
	}

	var result babyjub.PrivateKey
	copy(result[:], privateKey)

	return &result, nil
}
