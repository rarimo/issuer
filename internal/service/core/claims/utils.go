package claims

import (
	"crypto/rand"
	"encoding/binary"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CryptoRandUint64() (uint64, error) {
	var randBuffer [8]byte

	// TODO: We take only 4 bytes due to Polygon wallet issue
	if _, err := rand.Read(randBuffer[:4]); err != nil {
		return 0, errors.Wrap(err, "failed to read rand bytes")
	}
	return binary.LittleEndian.Uint64(randBuffer[:]), nil
}
