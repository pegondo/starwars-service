package swapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestComputePageIdxs(t *testing.T) {
	testCases := []struct {
		name               string
		offset             int
		remainingResources int
		apiPageSize        int
		idxs               indexes
	}{
		{
			name:               "zero_offset_and_remaining_resources_and_page_size",
			offset:             0,
			remainingResources: 0,
			apiPageSize:        0,
			idxs: indexes{
				min: 0,
				max: 0,
			},
		},
		{
			name:               "valid_range",
			offset:             2,
			remainingResources: 2,
			apiPageSize:        5,
			idxs: indexes{
				min: 2,
				max: 4,
			},
		},
		{
			name:               "too_many_remaining_resources",
			offset:             2,
			remainingResources: 10,
			apiPageSize:        5,
			idxs: indexes{
				min: 2,
				max: 5,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			idxs := computePageIdxs(tc.offset, tc.remainingResources, tc.apiPageSize)
			require.Equal(t, tc.idxs, idxs)
		})
	}
}

func TestComputeInitialPage(t *testing.T) {
	testCases := []struct {
		name        string
		pageNumber  int
		pageSize    int
		apiPageSize int
		page        page
	}{
		{
			name:        "one_page_number_and_page_size_and_api_page_size",
			pageNumber:  1,
			pageSize:    1,
			apiPageSize: 1,
			page: page{
				number: 1,
				offset: 0,
			},
		},
		{
			name:        "req_page_size_lower_than_api_page_size_first_page",
			pageNumber:  1,
			pageSize:    5,
			apiPageSize: 7,
			page: page{
				number: 1,
				offset: 0,
			},
		},
		{
			name:        "req_page_size_lower_than_api_page_size_second_page",
			pageNumber:  2,
			pageSize:    5,
			apiPageSize: 7,
			page: page{
				number: 1,
				offset: 5,
			},
		},
		{
			name:        "req_page_size_equal_than_api_page_size_first_page",
			pageNumber:  1,
			pageSize:    5,
			apiPageSize: 5,
			page: page{
				number: 1,
				offset: 0,
			},
		},
		{
			name:        "req_page_size_equal_than_api_page_size_second_page",
			pageNumber:  2,
			pageSize:    5,
			apiPageSize: 5,
			page: page{
				number: 2,
				offset: 0,
			},
		},
		{
			name:        "req_page_size_greater_than_api_page_size_first_page",
			pageNumber:  1,
			pageSize:    7,
			apiPageSize: 5,
			page: page{
				number: 1,
				offset: 0,
			},
		},
		{
			name:        "req_page_size_greater_than_api_page_size_second_page",
			pageNumber:  2,
			pageSize:    7,
			apiPageSize: 5,
			page: page{
				number: 2,
				offset: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			page := computeInitialPage(tc.pageNumber, tc.pageSize, tc.apiPageSize)
			require.Equal(t, tc.page, page)
		})
	}
}

func TestSortResults(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name            string
		results         []Person
		sortCriteria    SortCriteria
		expectedResults []Person
		err             error
	}{
		{
			name:            "invalid_sort_criteria",
			results:         nil,
			sortCriteria:    SortCriteria{},
			expectedResults: nil,
			err:             ErrInvalidSortField,
		},
		{
			name:    "nil_results",
			results: nil,
			sortCriteria: SortCriteria{
				Field: NameSortField,
			},
			expectedResults: nil,
			err:             nil,
		},
		{
			name:    "empty_results",
			results: []Person{},
			sortCriteria: SortCriteria{
				Field: NameSortField,
			},
			expectedResults: []Person{},
			err:             nil,
		},
		{
			name: "sort_by_name_asc",
			results: []Person{
				{
					Name: "1",
				},
				{
					Name: "3",
				},
				{
					Name: "2",
				},
			},
			sortCriteria: SortCriteria{
				Field: NameSortField,
				Order: AscendingOrder,
			},
			expectedResults: []Person{
				{
					Name: "1",
				},
				{
					Name: "2",
				},
				{
					Name: "3",
				},
			},
			err: nil,
		},
		{
			name: "sort_by_name_desc",
			results: []Person{
				{
					Name: "1",
				},
				{
					Name: "3",
				},
				{
					Name: "2",
				},
			},
			sortCriteria: SortCriteria{
				Field: NameSortField,
				Order: DescendingOrder,
			},
			expectedResults: []Person{
				{
					Name: "3",
				},
				{
					Name: "2",
				},
				{
					Name: "1",
				},
			},
			err: nil,
		},
		{
			name: "sort_by_created_asc",
			results: []Person{
				{
					Created: now.Add(2 * time.Hour),
				},
				{
					Created: now,
				},
				{
					Created: now.Add(time.Hour),
				},
			},
			sortCriteria: SortCriteria{
				Field: CreatedSortField,
				Order: AscendingOrder,
			},
			expectedResults: []Person{
				{
					Created: now,
				},
				{
					Created: now.Add(time.Hour),
				},
				{
					Created: now.Add(2 * time.Hour),
				},
			},
			err: nil,
		},
		{
			name: "sort_by_created_asc",
			results: []Person{
				{
					Created: now.Add(2 * time.Hour),
				},
				{
					Created: now,
				},
				{
					Created: now.Add(time.Hour),
				},
			},
			sortCriteria: SortCriteria{
				Field: CreatedSortField,
				Order: DescendingOrder,
			},
			expectedResults: []Person{
				{
					Created: now.Add(2 * time.Hour),
				},
				{
					Created: now.Add(time.Hour),
				},
				{
					Created: now,
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SortResults(tc.results, tc.sortCriteria)
			require.Equal(t, tc.expectedResults, tc.results)
			require.ErrorIs(t, tc.err, err)
		})
	}
}
