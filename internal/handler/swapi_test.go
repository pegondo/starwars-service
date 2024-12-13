package handler

import (
	"net/http"
	"testing"

	"github.com/pegondo/starwars/service/internal/resources/swapi"
	"github.com/stretchr/testify/require"
)

func TestGetStatusCode(t *testing.T) {
	testCases := []struct {
		name       string
		resp       swapi.SwapiResponse[swapi.Person]
		statusCode int
	}{
		{
			name: "zero_count_and_nil_results",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count:   0,
				Results: nil,
			},
			statusCode: http.StatusOK,
		},
		{
			name: "zero_count_and_empty_results",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count:   0,
				Results: []swapi.Person{},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "one_count_and_one_result",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count: 1,
				Results: []swapi.Person{
					{},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "two_count_and_two_results",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count: 2,
				Results: []swapi.Person{
					{}, {},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "five_count_and_five_results",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count: 5,
				Results: []swapi.Person{
					{}, {}, {}, {}, {},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "more_count_than_results",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count: 5,
				Results: []swapi.Person{
					{}, {},
				},
			},
			statusCode: http.StatusPartialContent,
		},
		{
			name: "more_results_than_count",
			resp: swapi.SwapiResponse[swapi.Person]{
				Count: 2,
				Results: []swapi.Person{
					{}, {}, {}, {}, {},
				},
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			statusCode := getStatusCode(tc.resp)
			require.Equal(t, tc.statusCode, statusCode)
		})
	}
}
