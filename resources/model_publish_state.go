/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PublishState struct {
	Key
	Attributes PublishStateAttributes `json:"attributes"`
}
type PublishStateResponse struct {
	Data     PublishState `json:"data"`
	Included Included     `json:"included"`
}

type PublishStateListResponse struct {
	Data     []PublishState `json:"data"`
	Included Included       `json:"included"`
	Links    *Links         `json:"links"`
}

// MustPublishState - returns PublishState from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPublishState(key Key) *PublishState {
	var publishState PublishState
	if c.tryFindEntry(key, &publishState) {
		return &publishState
	}
	return nil
}
