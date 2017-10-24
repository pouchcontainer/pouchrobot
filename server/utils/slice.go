package utils

// SliceContainsElement returns true if data is in current slice
func SliceContainsElement(input []string, data string) bool {
	for _, value := range input {
		if value == data {
			return true
		}
	}
	return false
}
