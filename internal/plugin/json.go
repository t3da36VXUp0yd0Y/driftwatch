package plugin

import (
	"encoding/json"

	"github.com/driftwatch/driftwatch/internal/drift"
)

// decodeJSON unmarshals a JSON byte slice into a drift.Result.
func decodeJSON(data []byte, r *drift.Result) error {
	return json.Unmarshal(data, r)
}
