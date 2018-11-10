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

import (
	"reflect"
	"testing"
)

func TestSliceContainsElement(t *testing.T) {
	type args struct {
		input []string
		data  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal contains",
			args: args{
				input: []string{"a", "b"},
				data:  "a",
			},
			want: true,
		},
		{
			name: "normal un-contains",
			args: args{
				input: []string{"a", "b"},
				data:  "aa",
			},
			want: false,
		},
		{
			name: "normal un-contains",
			args: args{
				input: []string{"a", "b"},
				data:  "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceContainsElement(tt.args.input, tt.args.data); got != tt.want {
				t.Errorf("SliceContainsElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueElementSlice(t *testing.T) {
	type args struct {
		data []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "already unique",
			args: args{
				data: []string{"a", "b"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "need to be unique",
			args: args{
				data: []string{"a", "a", "b"},
			},
			want: []string{"a", "b"},
		},
		{
			name: "nil slice",
			args: args{
				data: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqueElementSlice(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueElementSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
