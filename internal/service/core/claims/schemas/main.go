package schemas

import (
	"context"
	"encoding/json"
	"fmt"

	core "github.com/iden3/go-iden3-core"
	jsonSuite "github.com/iden3/go-schema-processor/json"
	"github.com/iden3/go-schema-processor/loaders"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/iden3/go-schema-processor/utils"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"

	"gitlab.com/q-dev/q-id/issuer/internal/service/core/claims"
)

func NewBuilder(ctx context.Context, schemasBaseUrl string) (*Builder, error) {
	builder := &Builder{
		SchemasBaseURL: schemasBaseUrl,
		CachedSchemas:  map[string]Schema{},
	}

	if err := builder.loadSchemas(ctx, schemasBaseUrl); err != nil {
		return nil, errors.Wrap(err, "failed to load schemas")
	}

	return builder, nil
}

func (b *Builder) loadSchemas(ctx context.Context, schemasBaseUrl string) error {
	for schemaType, schema := range claims.ClaimSchemaList {
		schemaBytes, _, err := (&loaders.HTTP{
			URL: fmt.Sprint(schemasBaseUrl, schema.ClaimSchemaURL),
		}).Load(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to load schema")
		}

		var parsedSchema jsonSuite.Schema
		if err = json.Unmarshal(schemaBytes, &parsedSchema); err != nil {
			return errors.Wrap(err, "failed to parse schema")
		}

		jsonLdContext, ok := parsedSchema.Metadata.Uris["jsonLdContext"].(string)
		if !ok {
			return errors.New("failed to get jsonLdContext from schema")
		}

		b.CachedSchemas[schemaType.ToRaw()] = Schema{
			Raw:           schemaBytes,
			Body:          parsedSchema,
			JSONLdContext: jsonLdContext,
		}
	}

	return nil
}

func (b *Builder) CreateCoreClaim(
	ctx context.Context,
	schemaType claims.ClaimSchemaType,
	credential *verifiable.W3CCredential,
	revNonce uint64,
) (*core.Claim, error) {
	parseOptions := &processor.CoreClaimOptions{
		RevNonce:        revNonce,
		Updatable:       false,
		SubjectPosition: claims.SubjectPositionIndex,
		MerklizedRootPosition: claims.DefineMerklizedRootPosition(
			b.CachedSchemas[schemaType.ToRaw()].Body.Metadata,
			utils.MerklizedRootPositionValue,
		),
	}

	claimsProcessor := processor.InitProcessorOptions(
		&processor.Processor{},
		processor.WithValidator(jsonSuite.Validator{}),
		processor.WithParser(jsonSuite.Parser{}),
	)

	jsonCredential, err := json.Marshal(credential)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal verifiable credential")
	}

	err = claimsProcessor.ValidateData(jsonCredential, b.CachedSchemas[schemaType.ToRaw()].Raw)
	if err != nil {
		return nil, errors.Wrap(ErrValidationData, err.Error())
	}

	coreClaim, err := claimsProcessor.ParseClaim(
		ctx,
		*credential,
		fmt.Sprintf("%s#%s", b.CachedSchemas[schemaType.ToRaw()].JSONLdContext, schemaType.ToRaw()),
		b.CachedSchemas[schemaType.ToRaw()].Raw,
		parseOptions,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse the schema slots")
	}

	return coreClaim, nil
}
