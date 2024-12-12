package request

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pegondo/starwars/service/internal/errors"
)

const (
	// pageParamKey is the key to get the page query parameter.
	pageParamKey = "page"
	// defaultPageParam is the default value for the page query parameter.
	defaultPageParam = 1

	// pageSizeParamKey is the key to get the page size query parameter.
	pageSizeParamKey = "pageSize"
	// defaultPageSize is the default page size used to request the API.
	defaultPageSize = 15

	// searchParamKey is the key to get the search query parameter.
	searchParamKey = "search"
	// defaultSearchValue is the default value for the search query parameter.
	defaultSearchValue = ""

	// sortFieldParamKey is the request parameter for the sort field.
	sortFieldParamKey = "sortField"
	// sortFieldParamKey is the request parameter for the sort order.
	sortOrderParamKey = "sortOrder"
)

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

// Validate validates if sc has valid data.
func (sc *SortCriteria) Validate() error {
	switch sc.Field {
	case "":
	case NameSortField:
	case CreatedSortField:
		// OK.
	default:
		return errors.New(errors.InvalidSortCriteriaErrorCode, errors.InvalidSortCriteriaErrorMsg)
	}
	switch sc.Order {
	case "":
	case AscendingOrder:
	case DescendingOrder:
		// OK.
	default:
		return errors.New(errors.InvalidSortCriteriaErrorCode, errors.InvalidSortCriteriaErrorMsg)
	}
	return nil
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

// RequestParams represents the parameters of the request.
type RequestParams struct {
	// Page is the number of the page requested.
	Page int
	// PageSize is the size of the page requested.
	PageSize int
	// Search is the search criteria requested.
	Search string
	// SortCriteria is the sorting criteria requested.
	SortCriteria *SortCriteria
}

// getNumericParam returns the parameter with the given key from the context. If
// the param isn't defined, getNumericParam returns the default value provided.
func getNumericParam(c *gin.Context, key string, defaultValue int) (value int, err error) {
	valueStr := c.DefaultQuery(key, strconv.Itoa(defaultValue))
	return strconv.Atoi(valueStr)
}

// Params extracts the request parameters from the context and returns them.
func Params(c *gin.Context) (params RequestParams, err error) {
	params = RequestParams{}

	params.Page, err = getNumericParam(c, pageParamKey, defaultPageParam)
	if err != nil || params.Page < 1 {
		return params, errors.New(errors.InvalidPageErrorCode, errors.InvalidPageErrorMsg)
	}

	params.PageSize, err = getNumericParam(c, pageSizeParamKey, defaultPageSize)
	if err != nil || params.PageSize < 1 {
		return params, errors.New(errors.InvalidPageSizeErrorCode, errors.InvalidPageSizeErrorMsg)
	}

	search := c.DefaultQuery(searchParamKey, defaultSearchValue)
	search = strings.ToLower(search)
	params.Search = search

	params.SortCriteria = GetSortCriteria(c)
	if params.SortCriteria != nil {
		if err = params.SortCriteria.Validate(); err != nil {
			return params, err
		}
	}

	return params, nil
}
