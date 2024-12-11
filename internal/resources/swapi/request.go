package swapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"starwars/service/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Source: https://swapi.dev/documentation
const (
	// swapiBaseUrl is the base URL of the SWAPI.
	swapiBaseUrl = "https://swapi.dev/api"
	// swapiPageSize is the SWAPI fixed page size.
	swapiPageSize = 10

	// sortFieldParamKey is the request parameter for the sort field.
	sortFieldParamKey = "sortField"
	// sortFieldParamKey is the request parameter for the sort order.
	sortOrderParamKey = "sortOrder"
)

// Resource represents a SWAPI resource the API serves.
type Resource interface {
	Person | Planet
	GetName() string
	GetCreated() time.Time
}

// SwapiResponse represents the SWAPI response for a resource T.
type SwapiResponse[T Resource] struct {
	// Count is the number of resources in the collection.
	Count int `json:"count"`
	// Next is the URL for the next page of the resource.
	Next *string `json:"next"`
	// Resoults are the paginated elements in the collection.
	Results []T `json:"results"`
}

// SortField represents a valid sort field.
type SortField string

const (
	// NameSortField represents a sorting by name.
	NameSortField SortField = "name"
	// CreatedSortField represents a sorting by creation date.
	CreatedSortField SortField = "created"
)

// SortOrder represents a valid sort order.
type SortOrder string

const (
	// AscendingOrder represents an ascending sorting order.
	AscendingOrder SortOrder = "asc"
	// DescendingOrder represents a descending sorting order.
	DescendingOrder SortOrder = "desc"
)

// SortCriteria represents a sort criteria.
type SortCriteria struct {
	// Field is the field to sort by.
	Field SortField
	// Order is the order to sort on.
	Order SortOrder
}

// GetSortCriteria returns the sorting criteria in the parameters of the given
// request context. If the request doesn't contain information about the
// sorting, GetSortCriteria returns nil.
func GetSortCriteria(c *gin.Context) *SortCriteria {
	sortField := c.DefaultQuery(sortFieldParamKey, "")
	sortField = strings.ToLower(sortField)
	if sortField == "" {
		return nil
	}

	sortOrder := c.DefaultQuery(sortOrderParamKey, string(AscendingOrder))
	sortOrder = strings.ToLower(sortOrder)

	return &SortCriteria{
		Field: SortField(sortField),
		Order: SortOrder(sortOrder),
	}
}

// request performs a HTTP request to the given URL returns its response.
func request[T Resource](url string) (response SwapiResponse[T], err error) {
	resp, err := http.Get(url)
	if err != nil {
		return response, fmt.Errorf("error while performing the request :: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error while reading the response body :: %v", err)
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, fmt.Errorf("error while parsing the response to a JSON :: %v", err)
	}

	return response, err
}

// retrievePageRec is a recursive solution to request SWAPI given a variable
// page size.
func retrievePageRec[T Resource](
	resource SwapiResponse[T],
	endpoint,
	search string,
	remainingResources,
	pageNumber,
	offset int,
) (
	swapiResp SwapiResponse[T],
	err error,
) {
	if remainingResources <= 0 {
		return resource, nil
	}

	url := fmt.Sprintf("%s/%s?page=%d", swapiBaseUrl, endpoint, pageNumber)
	if search != "" {
		url = fmt.Sprintf("%s&search=%s", url, search)
	}
	swapiResp, err = request[T](url)
	if err != nil {
		return swapiResp, fmt.Errorf("error while requesting the %s endpoint :: %v", endpoint, err)
	}
	if swapiResp.Count == 0 {
		// If there are no results, return an empty array.
		return SwapiResponse[T]{
			Count:   0,
			Results: []T{},
		}, nil
	}

	remainingResources = int(math.Min(float64(remainingResources), float64(swapiResp.Count)))

	minIdx := offset
	maxIdx := int(math.Min(swapiPageSize, float64(minIdx+remainingResources)))
	swapiResp.Results = append(resource.Results, swapiResp.Results[minIdx:maxIdx]...)
	numElementsAdded := maxIdx - minIdx
	return retrievePageRec(swapiResp, endpoint, search, remainingResources-numElementsAdded, pageNumber+1, 0)
}

// retrievePage retrieves the resources from the SWAPI with the given page
// number and size. If search isn't "", all the elements of resp.Results will
// contain the value of search.
func retrievePage[T Resource](
	endpoint string,
	page,
	pageSize int,
	search string,
) (
	resp SwapiResponse[T],
	err error,
) {
	numAlreadyRequestedResources := (page - 1) * pageSize
	initialPage := int(numAlreadyRequestedResources/swapiPageSize) + 1
	initialPageOffset := numAlreadyRequestedResources % swapiPageSize

	return retrievePageRec(SwapiResponse[T]{}, endpoint, search, pageSize, initialPage, initialPageOffset)
}

// retrieveAll returns all the resources in the given SWAPI endpoint. If search
// isn't "", the resources returned will contain the value of search in their
// name.
func retrieveAll[T Resource](
	endpoint,
	search string,
) (
	swapiResp SwapiResponse[T],
	err error,
) {
	url := fmt.Sprintf("%s/%s", swapiBaseUrl, endpoint)
	if search != "" {
		url = fmt.Sprintf("%s?search=%s", url, search)
	}

	swapiResp, err = request[T](url)
	if err != nil {
		return swapiResp, fmt.Errorf("error while requesting the %s endpoint :: %v", endpoint, err)
	}
	if swapiResp.Count == 0 {
		// If there are no results, return an empty array.
		return SwapiResponse[T]{
			Count:   0,
			Results: []T{},
		}, nil
	}

	for swapiResp.Next != nil {
		nextSwapiResp, err := request[T](*swapiResp.Next)
		if err != nil {
			return swapiResp, fmt.Errorf("error while requesting the %s endpoint :: %v", endpoint, err)
		}
		swapiResp.Results = append(swapiResp.Results, nextSwapiResp.Results...)
		swapiResp.Next = nextSwapiResp.Next
	}

	return swapiResp, err
}

// retrieveAllAndSort retrieves all the resources in SWAPI and sorts them using
// the given criteria to return the information paginated with the given page
// number and size. If search isn't "", the names of the resource in resp.Result
// will contain the value of search.
func retrieveAllAndSort[T Resource](
	endpoint string,
	page,
	pageSize int,
	search string,
	sortCriteria SortCriteria,
) (
	resp SwapiResponse[T],
	err error,
) {
	resources, err := retrieveAll[T](endpoint, search)
	if err != nil {
		return resp, err
	}

	var lessFn func(i, j int) bool
	switch sortCriteria.Field {
	case NameSortField:
		lessFn = func(i, j int) bool {
			return resources.Results[i].GetName() < resources.Results[j].GetName()
		}

	case CreatedSortField:
		lessFn = func(i, j int) bool {
			return resources.Results[i].GetCreated().Before(resources.Results[j].GetCreated())
		}
	default:
		return resp, fmt.Errorf("invalid sort field '%s'", sortCriteria.Field)
	}
	sort.Slice(resources.Results, lessFn)

	if sortCriteria.Order == DescendingOrder {
		utils.ReverseSlice(resources.Results)
	}

	minIdx := (page - 1) * pageSize
	if minIdx > len(resources.Results) {
		return SwapiResponse[T]{
			Count:   resources.Count,
			Results: []T{},
		}, nil
	}
	maxIdx := int(math.Min(float64(page*pageSize), float64(len(resources.Results))))
	resources.Results = resources.Results[minIdx:maxIdx]

	return resources, nil
}
