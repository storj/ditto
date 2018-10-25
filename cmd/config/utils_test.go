// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package config

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestContainsKey(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedResult bool
	}{
		{
			name:           "ValidCase",
			key:            "ListOptions.DefaultOptions.DefaultSource",
			expectedResult: true,
		},
		{
			name:           "Key is not valid",
			key:            "a",
			expectedResult: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contains := containsKey(test.key)

			assert.Equal(t, contains, test.expectedResult)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name           string
		item           string
		slice          []string
		expectedResult bool
	}{
		{
			name: "Empty slice",
			item: "a",
			slice: []string{

			},
			expectedResult: false,
		},
		{
			name: "Empty item",
			item: "",
			slice: []string{
				"a",
				"b",
			},
			expectedResult: false,
		},
		{
			name: "Valid case",
			item: "a",
			slice: []string{
				"a",
				"b",
			},
			expectedResult: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contains := contains(test.slice, test.item)

			assert.Equal(t, contains, test.expectedResult)
		})
	}
}
