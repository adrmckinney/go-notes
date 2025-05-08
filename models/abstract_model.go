package models

import "github.com/adrmckinney/go-notes/utils"

// FilterUpdateFields returns a new map containing only the key-value pairs from the input map
// whose keys are present in the allowed map. This is useful for sanitizing update payloads
// to ensure that only permitted fields are included in database such operations.
//
// Parameters:
//   - input:   map[string]interface{} containing all fields from the request or source.
//   - allowed: map[string]bool where keys are allowed field names and values are true. See Notes model for example
//
// Returns:
//   - map[string]interface{} containing only allowed fields from the input. Because GORM is receiving a map
//   - it assumes that all the column names are correct. Therefore the response converts to snake_case

func FilterUpdateFields(input map[string]any, allowed map[string]bool) map[string]any {
	filtered := make(map[string]any)
	for k, v := range input {
		if allowed[k] {
			filtered[utils.ToSnakeCase(k)] = v
		}
	}

	return filtered
}
