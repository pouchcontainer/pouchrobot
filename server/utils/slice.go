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

// DeltaSlice return a slice which contains element in compare but out of base.
func DeltaSlice(base, compare []string) []string {
	var result = []string{}
	for _, element := range compare {
		contained := false
		for _, baseElement := range base {
			if element == baseElement {
				contained = true
				continue
			}
		}
		if !contained {
			result = append(result, element)
		}
	}
	return result
}
