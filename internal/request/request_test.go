package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pegondo/starwars-service/internal/errors"
	"github.com/pegondo/starwars-service/internal/request"
	"github.com/stretchr/testify/require"
)

func buildRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.GET("/", handler)
	return r
}

func TestGetSortCriteria(t *testing.T) {
	testCases := []struct {
		name         string
		sortField    string
		sortOrder    string
		sortCriteria *request.SortCriteria
	}{
		{
			name:         "empty_sort_field_and_empty_sort_order",
			sortField:    "",
			sortOrder:    "",
			sortCriteria: nil,
		},
		{
			name:         "empty_sort_field_and_asc_sort_order",
			sortField:    "",
			sortOrder:    string(request.AscendingOrder),
			sortCriteria: nil,
		},
		{
			name:         "empty_sort_field_and_desc_sort_order",
			sortField:    "",
			sortOrder:    string(request.DescendingOrder),
			sortCriteria: nil,
		},
		{
			name:      "name_sort_field_and_empty_sort_order",
			sortField: string(request.NameSortField),
			sortOrder: "",
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: "",
			},
		},
		{
			name:      "created_sort_field_and_empty_sort_order",
			sortField: string(request.CreatedSortField),
			sortOrder: "",
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: "",
			},
		},
		{
			name:      "name_sort_field_and_asc_sort_order",
			sortField: string(request.NameSortField),
			sortOrder: string(request.AscendingOrder),
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: request.AscendingOrder,
			},
		},
		{
			name:      "name_sort_field_and_desc_sort_order",
			sortField: string(request.NameSortField),
			sortOrder: string(request.DescendingOrder),
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: request.DescendingOrder,
			},
		},
		{
			name:      "created_sort_field_and_asc_sort_order",
			sortField: string(request.CreatedSortField),
			sortOrder: string(request.AscendingOrder),
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: request.AscendingOrder,
			},
		},
		{
			name:      "created_sort_field_and_desc_sort_order",
			sortField: string(request.CreatedSortField),
			sortOrder: string(request.DescendingOrder),
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: request.DescendingOrder,
			},
		},
		{
			name:      "invalid_sort_field",
			sortField: "<invalid-sort-field>",
			sortOrder: "",
			sortCriteria: &request.SortCriteria{
				Field: "<invalid-sort-field>",
				Order: "",
			},
		},
		{
			name:         "invalid_sort_order_empty_sort_field",
			sortField:    "",
			sortOrder:    "<invalid-sort-order>",
			sortCriteria: nil,
		},
		{
			name:      "invalid_sort_order_name_sort_field",
			sortField: string(request.NameSortField),
			sortOrder: "<invalid-sort-order>",
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: "<invalid-sort-order>",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var sortCriteria *request.SortCriteria
			handler := func(c *gin.Context) {
				sortCriteria = request.GetSortCriteria(c)
			}
			r := buildRouter(handler)

			url := fmt.Sprintf("/?sortField=%s&sortOrder=%s", tc.sortField, tc.sortOrder)
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.sortCriteria, sortCriteria)
		})
	}
}

func TestValidateSortCriteria(t *testing.T) {
	testCases := []struct {
		name         string
		sortCriteria *request.SortCriteria
		err          error
	}{
		{
			name: "empty_sort_field_and_asc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: "",
				Order: request.AscendingOrder,
			},
			err: nil,
		},
		{
			name: "empty_sort_field_and_desc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: "",
				Order: request.DescendingOrder,
			},
			err: nil,
		},
		{
			name: "name_sort_field_and_empty_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: "",
			},
			err: nil,
		},
		{
			name: "created_sort_field_and_empty_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: "",
			},
			err: nil,
		},
		{
			name: "name_sort_field_and_asc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: request.AscendingOrder,
			},
			err: nil,
		},
		{
			name: "name_sort_field_and_desc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.NameSortField,
				Order: request.DescendingOrder,
			},
			err: nil,
		},
		{
			name: "created_sort_field_and_asc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: request.AscendingOrder,
			},
			err: nil,
		},
		{
			name: "created_sort_field_and_desc_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: request.CreatedSortField,
				Order: request.DescendingOrder,
			},
			err: nil,
		},
		{
			name: "invalid_sort_field",
			sortCriteria: &request.SortCriteria{
				Field: "<invalid-sort-field>",
				Order: "",
			},
			err: errors.New(errors.InvalidSortCriteriaErrorCode, errors.InvalidSortCriteriaErrorMsg),
		},
		{
			name: "invalid_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: "",
				Order: "<invalid-sort-order>",
			},
			err: errors.New(errors.InvalidSortCriteriaErrorCode, errors.InvalidSortCriteriaErrorMsg),
		},
		{
			name: "invalid_sort_field_and_invalid_sort_order",
			sortCriteria: &request.SortCriteria{
				Field: "<invalid-sort-field>",
				Order: "<invalid-sort-order>",
			},
			err: errors.New(errors.InvalidSortCriteriaErrorCode, errors.InvalidSortCriteriaErrorMsg),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sortCriteria.Validate()
			require.Equal(t, tc.err, err)
		})
	}
}

func TestGetParams(t *testing.T) {
	testCases := []struct {
		name      string
		page      string
		pageSize  string
		search    string
		sortField string
		sortOrder string
		params    request.RequestParams
		err       error
	}{
		// Page number.
		{
			name:     "emtpy_page",
			page:     "",
			pageSize: "1",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg),
		},
		{
			name:     "valid_page",
			page:     "3",
			pageSize: "1",
			params: request.RequestParams{
				Page:     3,
				PageSize: 1,
			},
			err: nil,
		},
		{
			name:     "zero_page",
			page:     "0",
			pageSize: "1",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg),
		},
		{
			name:     "negative_page",
			page:     "-1",
			pageSize: "1",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg),
		},
		{
			name:     "invalid_page",
			page:     "<invalid-page>",
			pageSize: "1",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg),
		},
		// Page size.
		{
			name:     "emtpy_page_size",
			page:     "1",
			pageSize: "",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg),
		},
		{
			name:     "valid_page_size",
			page:     "1",
			pageSize: "2",
			params: request.RequestParams{
				Page:     1,
				PageSize: 2,
			},
			err: nil,
		},
		{
			name:     "zero_page_size",
			page:     "1",
			pageSize: "0",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg),
		},
		{
			name:     "negative_page_size",
			page:     "1",
			pageSize: "-1",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg),
		},
		{
			name:     "invalid_page_size",
			page:     "1",
			pageSize: "<invalid-page-size>",
			params:   request.RequestParams{},
			err:      errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg),
		},
		{
			name:     "empty_search",
			page:     "1",
			pageSize: "1",
			search:   "",
			params: request.RequestParams{
				Page:     1,
				PageSize: 1,
				Search:   "",
			},
			err: nil,
		},
		{
			name:     "valid_search",
			page:     "1",
			pageSize: "1",
			search:   "<search>",
			params: request.RequestParams{
				Page:     1,
				PageSize: 1,
				Search:   "<search>",
			},
			err: nil,
		},
		{
			name:     "capitalized_search",
			page:     "1",
			pageSize: "1",
			search:   "<SeArCH>",
			params: request.RequestParams{
				Page:     1,
				PageSize: 1,
				Search:   "<search>",
			},
			err: nil,
		},
		// The sort criteria is already tested in TestValidateSortCriteria.

		// All the params.
		{
			name:      "all_the_params",
			page:      "1",
			pageSize:  "1",
			search:    "<search>",
			sortField: string(request.NameSortField),
			sortOrder: string(request.AscendingOrder),
			params: request.RequestParams{
				Page:     1,
				PageSize: 1,
				Search:   "<search>",
				SortCriteria: &request.SortCriteria{
					Field: request.NameSortField,
					Order: request.AscendingOrder,
				},
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var params request.RequestParams
			var err error
			handler := func(c *gin.Context) {
				params, err = request.Params(c)
			}
			r := buildRouter(handler)

			url := fmt.Sprintf("/?page=%s&pageSize=%s&search=%s&sortField=%s&sortOrder=%s",
				tc.page, tc.pageSize, tc.search, tc.sortField, tc.sortOrder)
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.params, params)
			}
		})
	}
}
