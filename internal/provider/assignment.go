package provider

import (
	"slices"
)

// assignKeys assigns counter values to the keys provided as input
func assignKeys(keys []string, state map[string]int64, reuse bool, initial int64, last int64) (int64, map[string]int64) {
	// Create a map to hold the assigned values
	assignedValues := make(map[string]int64, len(keys))
	// Create a list of values for easier tracking for the next possible one
	values := make([]int64, 0, len(keys))

	// If the previous state is defined, maintain all entries that are still present in keys.
	// Also handle a changing initial value.
	for key, value := range state {
		if slices.Contains(keys, key) && value >= initial {
			assignedValues[key] = value
			values = append(values, value)
		}
	}

	// Sort keys to provide a predictable behaviour
	slices.Sort(keys)

	// Iterate over the keys and provide values to those not covered yet
	for _, key := range keys {
		// If the key has not yet a value assigned
		if _, exists := assignedValues[key]; !exists {
			// If reuse is true, find a value that does not exist in the assignedValues map
			if reuse {
				for i := initial; ; i++ {
					if !slices.Contains(values, i) {
						assignedValues[key] = i
						values = append(values, i)
						last = i
						break
					}
				}
			} else {
				// If reuse is false, increment the last value and assign it to the key
				last = max(last+1, initial)
				assignedValues[key] = last
			}
		}
	}

	// Return the last value and the assignedValues map
	return last, assignedValues
}
