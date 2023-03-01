/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type IssueClaimAttributes struct {
	Credential json.RawMessage `json:"credential"`
	// The claim expiration date in RFC3339 format
	Expiration string `json:"expiration"`
}
