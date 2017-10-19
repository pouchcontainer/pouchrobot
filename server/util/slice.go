package util

// uniqueSlice returns a slice without two same elements
func uniqueSlice(input []string) []string {
	if input == nil {
		return nil
	}

	result := make([]string, 0, len(input))
	internal := make(map[string]struct{}, len(input))

	for _, value := range input {
		if _, exist := internal[value]; !exist {
			internal[value] = struct{}{}
		}
	}

	for key := range internal {
		result = append(result, key)
	}
	return result
}

// SliceContainsElement judges whether an element is includes in an slice
func SliceContainsElement(origin []string, input string) bool {
	for _, value := range origin {
		if value == input {
			return true
		}
	}
	return false
}
