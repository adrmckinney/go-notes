package models

// FilterUpdateFields returns a new map containing only the key-value pairs from the input map
// whose keys are present in the allowed map. This is useful for sanitizing update payloads
// to ensure that only permitted fields are included in database such operations.
//
// Parameters:
//   - input:   map[string]interface{} containing all fields from the request or source.
//   - allowed: map[string]bool where keys are allowed field names and values are true. See Notes model for example
//
// Returns:
//   - map[string]interface{} containing only allowed fields from the input.
func FilterUpdateFields(input map[string]interface{}, allowed map[string]bool) map[string]interface{} {
	filtered := make(map[string]interface{})
	for k, v := range input {
		if allowed[k] {
			filtered[k] = v
		}
	}
	return filtered
}
