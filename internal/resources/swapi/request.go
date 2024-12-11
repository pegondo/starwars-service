package swapi

import (
	"encoding/json"
	"errors"
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

// ErrInvalidSortField is the error returned when the sort field in SortCriteria
// is invalid.
var ErrInvalidSortField = errors.New("invalid sort field")

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

// indexes represent a pair of min and max indexes.
type indexes struct {
	min int
	max int
}

// page represents the information needed to correctly fetch a SWAPI page.
type page struct {
	number int
	offset int
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

// buildUrl builds the SWAPI URL to request with the given endpoint, page number
// and search condition.
func buildUrl(endpoint string, pageNumber int, search string) string {
	url := fmt.Sprintf("%s/%s?page=%d", swapiBaseUrl, endpoint, pageNumber)
	if search != "" {
		url = fmt.Sprintf("%s&search=%s", url, search)
	}
	return url
}

// computePageIdxs returns the min and max page indexes needed to correctly
// request for the page with the given offset and remaining resources.
// computePageIdxs will misbehave if a negative offset or page size is provided,
// or if the offset is bigger than the page size.
func computePageIdxs(offset int, remainingResources, apiPageSize int) indexes {
	min := offset
	max := int(math.Min(float64(apiPageSize), float64(min+remainingResources)))
	return indexes{
		min: min,
		max: max,
	}
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

	url := buildUrl(endpoint, pageNumber, search)
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

	idxs := computePageIdxs(offset, remainingResources, swapiPageSize)
	numElementsAdded := idxs.max - idxs.min
	swapiResp.Results = append(resource.Results, swapiResp.Results[idxs.min:idxs.max]...)

	return retrievePageRec(swapiResp, endpoint, search, remainingResources-numElementsAdded, pageNumber+1, 0)
}

// computeInitialPage computes the number of the page to request and its offset
// based on the given page number, requested page size and API page size. If the
// page number, page size or API page size are lower than one,
// computeInitialPage may misbehave.
func computeInitialPage(pageNumber, pageSize, apiPageSize int) page {
	numAlreadyRequestedResources := (pageNumber - 1) * pageSize
	initial := int(numAlreadyRequestedResources/apiPageSize) + 1
	offset := numAlreadyRequestedResources % apiPageSize
	return page{
		number: initial,
		offset: offset,
	}
}

// retrievePage retrieves the resources from the SWAPI with the given page
// number and size. If search isn't "", all the elements of resp.Results will
// contain the value of search.
func retrievePage[T Resource](
	endpoint string,
	pageNumber,
	pageSize int,
	search string,
) (
	resp SwapiResponse[T],
	err error,
) {
	page := computeInitialPage(pageNumber, pageSize, swapiPageSize)
	return retrievePageRec(SwapiResponse[T]{}, endpoint, search, pageSize, page.number, page.offset)
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
	url := buildUrl(endpoint, 1, search)

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

// sortResults sorts the given slice of results based on the given sort
// criteria.
func sortResults[T Resource](results []T, sortCriteria SortCriteria) error {
	var lessFn func(i, j int) bool
	switch sortCriteria.Field {
	case NameSortField:
		lessFn = func(i, j int) bool {
			return results[i].GetName() < results[j].GetName()
		}

	case CreatedSortField:
		lessFn = func(i, j int) bool {
			return results[i].GetCreated().Before(results[j].GetCreated())
		}
	default:
		return ErrInvalidSortField
	}
	sort.Slice(results, lessFn)

	if sortCriteria.Order == DescendingOrder {
		utils.ReverseSlice(results)
	}
	return nil
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

	if err = sortResults(resources.Results, sortCriteria); err != nil {
		return resp, err
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
