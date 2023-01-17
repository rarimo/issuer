package zkp

import (
	"bytes"
	"math/big"
	"strings"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *ZKProof) ProofToBigInts() (a []*big.Int, b [][]*big.Int, c []*big.Int, err error) {

	a, err = arrayStringToBigInt(p.A)
	if err != nil {
		return nil, nil, nil, err
	}
	b = make([][]*big.Int, len(p.B))
	for i, v := range p.B {
		b[i], err = arrayStringToBigInt(v)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	c, err = arrayStringToBigInt(p.C)
	if err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

func arrayStringToBigInt(str []string) ([]*big.Int, error) {
	var result []*big.Int
	for i := 0; i < len(str); i++ {
		si, err := stringToBigInt(str[i])
		if err != nil {
			return nil, err
		}
		result = append(result, si)
	}
	return result, nil
}

func stringToBigInt(str string) (*big.Int, error) {
	base := 10
	if bytes.HasPrefix([]byte(str), []byte("0x")) {
		base = 16
		str = strings.TrimPrefix(str, "0x")
	}
	n, ok := new(big.Int).SetString(str, base)
	if !ok {
		return nil, errors.Errorf("can not parse string to *big.Int: %s", str)
	}
	return n, nil
}
