package claims

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"

	"github.com/iancoleman/strcase"
	"github.com/iden3/go-schema-processor/verifiable"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CryptoRandUint64() (uint64, error) {
	var randBuffer [8]byte

	// TODO: We take only 4 bytes due to PolygonID wallet issue
	if _, err := rand.Read(randBuffer[:4]); err != nil {
		return 0, errors.Wrap(err, "failed to read rand bytes")
	}
	return binary.LittleEndian.Uint64(randBuffer[:]), nil
}

func ParseCredentialFromSnakeCase(credentialRaw []byte) (*verifiable.W3CCredential, error) {
	credentialsMap := map[string]interface{}{}
	if err := json.Unmarshal(credentialRaw, &credentialsMap); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal credentials to map")
	}

	credentials := verifiable.W3CCredential{
		CredentialSubject: map[string]interface{}{},
	}
	for key, value := range credentialsMap {
		credentials.CredentialSubject[strcase.ToLowerCamel(key)] = value
	}

	return &credentials, nil
}
