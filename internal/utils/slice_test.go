package utils_test

import (
	"testing"

	"github.com/pegondo/starwars-service/internal/utils"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	name string
}

func TestRevertSlice(t *testing.T) {
	testCases := []struct {
		name          string
		slice         []testStruct
		expectedSlice []testStruct
	}{
		{
			name:          "nil_slice",
			slice:         nil,
			expectedSlice: nil,
		},
		{
			name:          "empty_slice",
			slice:         []testStruct{},
			expectedSlice: []testStruct{},
		},
		{
			name: "single_element_slice",
			slice: []testStruct{
				{
					name: "<element-1>",
				},
			},
			expectedSlice: []testStruct{
				{
					name: "<element-1>",
				},
			},
		},
		{
			name: "two_elements_slice",
			slice: []testStruct{
				{
					name: "<element-1>",
				},
				{
					name: "<element-2>",
				},
			},
			expectedSlice: []testStruct{
				{
					name: "<element-2>",
				},
				{
					name: "<element-1>",
				},
			},
		},
		{
			name: "many_elements_slice",
			slice: []testStruct{
				{
					name: "<element-1>",
				},
				{
					name: "<element-2>",
				},
				{
					name: "<element-3>",
				},
				{
					name: "<element-4>",
				},
				{
					name: "<element-5>",
				},
			},
			expectedSlice: []testStruct{
				{
					name: "<element-5>",
				},
				{
					name: "<element-4>",
				},
				{
					name: "<element-3>",
				},
				{
					name: "<element-2>",
				},
				{
					name: "<element-1>",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			utils.ReverseSlice(tc.slice)
			require.Equal(t, tc.expectedSlice, tc.slice)
		})
	}
}
