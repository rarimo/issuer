package schemas

import (
	"context"
	"encoding/hex"
	"net/url"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-schema-processor/json"
	jsonld "github.com/iden3/go-schema-processor/json-ld"
	"github.com/iden3/go-schema-processor/loaders"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/pkg/errors"
)

func NewBuilder(ipfsURL string) (*Builder, error) {
	parsedIPFSURL, err := url.Parse(ipfsURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ipfs url")
	}

	return &Builder{
		ipfsURL: parsedIPFSURL,
	}, nil
}

func (b *Builder) Process(
	ctx context.Context,
	data []byte,
	schemaType string,
	schemaURLRaw string,
) (*processor.ParsedSlots, string, error) {
	schemaURL, err := url.Parse(schemaURLRaw)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse schema url")
	}

	slots, schemaBytes, err := b.getParsedSlots(ctx, schemaURL, schemaType, data)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to get parsed slots from schema data")
	}

	return slots, b.createSchemaHash(schemaBytes, schemaType), nil
}

// nolint
func (b *Builder) load(ctx context.Context, schemaURL *url.URL) (schema []byte, err error) {
	loader, err := b.getLoader(schemaURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get schema loader")
	}

	var schemaBytes []byte
	schemaBytes, _, err = loader.Load(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load schema")
	}

	return schemaBytes, nil
}

func (b *Builder) getParsedSlots(
	ctx context.Context,
	schemaURL *url.URL,
	schemaType string,
	dataBytes []byte,
) (*processor.ParsedSlots, []byte, error) {
	loader, err := b.getLoader(schemaURL)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get schema loader")
	}

	schema, schemaFormat, err := loader.Load(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to load the schema")
	}

	schemaProcessor, err := newSchemaProcessor(schemaType, schemaFormat, loader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create new schema processor")
	}

	err = schemaProcessor.ValidateData(dataBytes, schema)
	if err != nil {
		return nil, nil, errors.Wrap(ErrValidationData, err.Error())
	}

	schemaSlots, err := schemaProcessor.ParseSlots(dataBytes, schema)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse the schema slots")
	}

	return &schemaSlots, schema, nil
}

func (b *Builder) getLoader(schemaURL *url.URL) (processor.SchemaLoader, error) {
	switch schemaURL.Scheme {
	case httpProtocolName, httpsProtocolName:
		return &loaders.HTTP{
			URL: schemaURL.String(),
		}, nil
	case ipfsProtocolName:
		return loaders.IPFS{
			URL: schemaURL.String(),
			CID: schemaURL.Host,
		}, nil
	default:
		return nil, errors.New("protocol is invalid or not supported")
	}
}

func newSchemaProcessor(
	schemaType string,
	schemaFormat string,
	loader processor.SchemaLoader,
) (*processor.Processor, error) {
	switch schemaFormat {
	case SchemaFormatJSON:
		return newJSONSchemaProcessor(schemaType, loader), nil
	case SchemaFormatJSONLD:
		return newJSONldSchemaProcessor(schemaType, loader), nil
	default:
		return nil, ErrSchemaFormatIsNotSupported
	}
}

func newJSONldSchemaProcessor(schemaType string, loader processor.SchemaLoader) *processor.Processor {
	return processor.InitProcessorOptions(
		&processor.Processor{},
		processor.WithValidator(jsonld.Validator{
			ClaimType: schemaType,
		}),
		processor.WithParser(jsonld.Parser{
			ClaimType:       schemaType,
			ParsingStrategy: processor.SlotFullfilmentStrategy,
		}),
		processor.WithSchemaLoader(loader),
	)
}

func newJSONSchemaProcessor(schemaType string, loader processor.SchemaLoader) *processor.Processor {
	return processor.InitProcessorOptions(
		&processor.Processor{},
		processor.WithValidator(json.Validator{}),
		processor.WithParser(jsonld.Parser{
			ClaimType:       schemaType,
			ParsingStrategy: processor.SlotFullfilmentStrategy,
		}),
		processor.WithSchemaLoader(loader),
	)
}

func (b *Builder) createSchemaHash(schemaBytes []byte, credentialType string) string {
	schemaHash := crypto.Keccak256(schemaBytes, []byte(credentialType))
	return hex.EncodeToString(schemaHash[len(schemaHash)-16:])
}
