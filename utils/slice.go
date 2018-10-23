// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if data == nil {
		return nil
	}

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

// DeltaSlice returns a slice which contains element in compare but out of base.
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

// SliceContainsSlice returns true if a slice contains another one
func SliceContainsSlice(old, new []string) bool {
	for _, newElement := range new {
		in := false
		for _, oldElement := range old {
			if newElement == oldElement {
				in = true
			}
		}
		if !in {
			return false
		}
	}
	return true
}
