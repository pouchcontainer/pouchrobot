package utils

// SliceContainsElement returns true if data is in current slice.
func SliceContainsElement(input []string, data string) bool {
	for _, value := range input {
		if value == data {
			return true
		}
	}
	return false
}

// UniqueElementSlice returns a slice with no element duplicated.
func UniqueElementSlice(data []string) []string {
	dataMap := make(map[string]struct{}, len(data))
	for _, value := range data {
		if _, exist := dataMap[value]; !exist {
			dataMap[value] = struct{}{}
		}
	}
	noDuplicatedSlice := make([]string, 0, len(dataMap))
	for key := range dataMap {
		noDuplicatedSlice = append(noDuplicatedSlice, key)
	}
	return noDuplicatedSlice

}
