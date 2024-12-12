package tests

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/pegondo/starwars/service/ex/client"
	"github.com/pegondo/starwars/service/internal/errors"
	"github.com/pegondo/starwars/service/internal/request"
	"github.com/pegondo/starwars/service/internal/resources/swapi"
	"github.com/stretchr/testify/require"
)

func TestRetrievePlanets_Page1(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, pageSize, len(resp.Response.Data))
}

func TestRetrievePlanets_Page2(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("2", pageSizeStr, "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, pageSize, len(resp.Response.Data))
}

func TestRetrievePlanets_PageOutOfIndex(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts(strconv.Itoa(math.MaxInt), pageSizeStr, "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.Empty(t, resp.Response.Data)
}

func TestRetrievePlanets_Page1_PageSize1(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", "1", "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 1, len(resp.Response.Data))
}

func TestRetrievePlanets_Page2_PageSize1(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("2", "1", "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 1, len(resp.Response.Data))
}

func TestRetrievePlanets_Page1_PageSize2(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", "2", "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 2, len(resp.Response.Data))
}

func TestRetrievePlanets_Page2_PageSize2(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("2", "2", "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 2, len(resp.Response.Data))
}

func TestRetrievePlanets_Page1_PageSizeTooBig(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", strconv.Itoa(math.MaxInt), "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
}

func TestRetrievePlanets_Page2_PageSizeTooBig(t *testing.T) {
	resp, err := c.RetrievePlanets(client.NewRequestOpts("2", strconv.Itoa(math.MaxInt), "", "", ""))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.Empty(t, resp.Response.Data)
}

func TestRetrievePlanets_Page1_Search(t *testing.T) {
	search := "a"
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, search, "", ""))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that all the planets name contain the searched segment.
	for _, planet := range resp.Response.Data {
		require.Contains(t, strings.ToLower(planet.Name), strings.ToLower(search))
	}
}

func TestRetrievePlanets_Page1_SortByNameAsc(t *testing.T) {
	sortCriteria := request.SortCriteria{
		Field: request.NameSortField,
		Order: request.AscendingOrder,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", string(sortCriteria.Field), string(sortCriteria.Order)))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the planets are sorted by name in ascending order.
	sortedPlanets := make([]swapi.Planet, len(resp.Response.Data))
	copy(sortedPlanets, resp.Response.Data)
	swapi.SortResults(sortedPlanets, sortCriteria)
	require.Equal(t, sortedPlanets, resp.Response.Data)
}

func TestRetrievePlanets_Page1_SortByNameDesc(t *testing.T) {
	sortCriteria := request.SortCriteria{
		Field: request.NameSortField,
		Order: request.DescendingOrder,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", string(sortCriteria.Field), string(sortCriteria.Order)))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the planets are sorted by name in descending order.
	sortedPlanets := make([]swapi.Planet, len(resp.Response.Data))
	copy(sortedPlanets, resp.Response.Data)
	swapi.SortResults(sortedPlanets, sortCriteria)
	require.Equal(t, sortedPlanets, resp.Response.Data)
}

func TestRetrievePlanets_Page1_SortByCreatedAsc(t *testing.T) {
	sortCriteria := request.SortCriteria{
		Field: request.CreatedSortField,
		Order: request.AscendingOrder,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", string(sortCriteria.Field), string(sortCriteria.Order)))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the planets are sorted by creation date in ascending order.
	sortedPlanets := make([]swapi.Planet, len(resp.Response.Data))
	copy(sortedPlanets, resp.Response.Data)
	swapi.SortResults(sortedPlanets, sortCriteria)
	require.Equal(t, sortedPlanets, resp.Response.Data)
}

func TestRetrievePlanets_Page1_SortByCreatedDesc(t *testing.T) {
	sortCriteria := request.SortCriteria{
		Field: request.CreatedSortField,
		Order: request.DescendingOrder,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", string(sortCriteria.Field), string(sortCriteria.Order)))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the planets are sorted by creation date in descending order.
	sortedPlanets := make([]swapi.Planet, len(resp.Response.Data))
	copy(sortedPlanets, resp.Response.Data)
	swapi.SortResults(sortedPlanets, sortCriteria)
	require.Equal(t, sortedPlanets, resp.Response.Data)
}

func TestRetrievePlanets_InvalidPageParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageErrorCode,
		ErrorMessage: errors.InvalidPageErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("<invalid-page>", pageSizeStr, "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_ZeroPageParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageErrorCode,
		ErrorMessage: errors.InvalidPageErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("0", pageSizeStr, "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_NegativePageParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageErrorCode,
		ErrorMessage: errors.InvalidPageErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("-1", pageSizeStr, "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_InvalidPageSizeParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageSizeErrorCode,
		ErrorMessage: errors.InvalidPageSizeErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", "<invalid-page-size>", "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_ZeroPageSizeParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageSizeErrorCode,
		ErrorMessage: errors.InvalidPageSizeErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", "0", "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_NegativePageSizeParam(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidPageSizeErrorCode,
		ErrorMessage: errors.InvalidPageSizeErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", "-1", "", "", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}

func TestRetrievePlanets_InvalidSortField(t *testing.T) {
	expectedErr := &errors.ResponseError{
		ErrorCode:    errors.InvalidSortCriteriaErrorCode,
		ErrorMessage: errors.InvalidSortCriteriaErrorMsg,
	}
	resp, err := c.RetrievePlanets(client.NewRequestOpts("1", pageSizeStr, "", "<invalid-sort-field>", ""))

	require.Equal(t, expectedErr, err)
	require.Nil(t, resp.Response.Data)
}
