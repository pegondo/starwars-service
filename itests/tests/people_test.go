package tests

import (
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/pegondo/starwars/service/ex/client"
	"github.com/pegondo/starwars/service/internal/resources/swapi"
	"github.com/stretchr/testify/require"
)

func TestRetrievePeople_Page1(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, pageSize, len(resp.Response.Data))
}

func TestRetrievePeople_Page2(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(2, pageSize, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, pageSize, len(resp.Response.Data))
}

func TestRetrievePeople_PageOutOfIndex(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(math.MaxInt, pageSize, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.Empty(t, resp.Response.Data)
}

func TestRetrievePeople_Page1_PageSize1(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, 1, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 1, len(resp.Response.Data))
}

func TestRetrievePeople_Page2_PageSize1(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(2, 1, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 1, len(resp.Response.Data))
}

func TestRetrievePeople_Page1_PageSize2(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, 2, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 2, len(resp.Response.Data))
}

func TestRetrievePeople_Page2_PageSize2(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(2, 2, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusPartialContent, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
	require.Equal(t, 2, len(resp.Response.Data))
}

func TestRetrievePeople_Page1_PageSizeTooBig(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, math.MaxInt, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)
}

func TestRetrievePeople_Page2_PageSizeTooBig(t *testing.T) {
	resp, err := c.RetrievePeople(client.NewRequestOpts(2, math.MaxInt, "", nil))

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.Empty(t, resp.Response.Data)
}

func TestRetrievePeople_Page1_Search(t *testing.T) {
	search := "sky"
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, search, nil))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that all the people name contain the searched segment.
	for _, person := range resp.Response.Data {
		require.Contains(t, strings.ToLower(person.Name), strings.ToLower(search))
	}
}

func TestRetrievePeople_Page1_SortByNameAsc(t *testing.T) {
	sortCriteria := swapi.SortCriteria{
		Field: swapi.NameSortField,
		Order: swapi.AscendingOrder,
	}
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, "", &sortCriteria))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the people are sorted by name in ascending order.
	sortedPeople := make([]swapi.Person, len(resp.Response.Data))
	copy(sortedPeople, resp.Response.Data)
	swapi.SortResults(sortedPeople, sortCriteria)
	require.Equal(t, sortedPeople, resp.Response.Data)
}

func TestRetrievePeople_Page1_SortByNameDesc(t *testing.T) {
	sortCriteria := swapi.SortCriteria{
		Field: swapi.NameSortField,
		Order: swapi.DescendingOrder,
	}
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, "", &sortCriteria))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the people are sorted by name in descending order.
	sortedPeople := make([]swapi.Person, len(resp.Response.Data))
	copy(sortedPeople, resp.Response.Data)
	swapi.SortResults(sortedPeople, sortCriteria)
	require.Equal(t, sortedPeople, resp.Response.Data)
}

func TestRetrievePeople_Page1_SortByCreatedAsc(t *testing.T) {
	sortCriteria := swapi.SortCriteria{
		Field: swapi.CreatedSortField,
		Order: swapi.AscendingOrder,
	}
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, "", &sortCriteria))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the people are sorted by creation date in ascending order.
	sortedPeople := make([]swapi.Person, len(resp.Response.Data))
	copy(sortedPeople, resp.Response.Data)
	swapi.SortResults(sortedPeople, sortCriteria)
	require.Equal(t, sortedPeople, resp.Response.Data)
}

func TestRetrievePeople_Page1_SortByCreatedDesc(t *testing.T) {
	sortCriteria := swapi.SortCriteria{
		Field: swapi.CreatedSortField,
		Order: swapi.DescendingOrder,
	}
	resp, err := c.RetrievePeople(client.NewRequestOpts(1, pageSize, "", &sortCriteria))

	require.NoError(t, err)
	require.Contains(t, []int{http.StatusOK, http.StatusPartialContent}, resp.StatusCode)
	require.NotNil(t, resp.Response.Data)
	require.NotEmpty(t, resp.Response.Data)

	// Verify that the people are sorted by creation date in descending order.
	sortedPeople := make([]swapi.Person, len(resp.Response.Data))
	copy(sortedPeople, resp.Response.Data)
	swapi.SortResults(sortedPeople, sortCriteria)
	require.Equal(t, sortedPeople, resp.Response.Data)
}
