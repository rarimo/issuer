/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type IssueClaimAttributes struct {
	// The claim expiration date in RFC3339 format
	Expiration string          `json:"expiration"`
	SchemaData json.RawMessage `json:"schema_data"`
	// The schema type
	SchemaType string `json:"schema_type"`
}
