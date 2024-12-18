package errors

const (
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
	InternalServerErrorMsg  = "An internal server error occurred."

	InvalidPageErrorCode = "INVALID_PAGE"
	InvalidPageErrorMsg  = "The page must be a number greater than 0."

	InvalidPageSizeErrorCode = "INVALID_PAGE_SIZE"
	InvalidPageSizeErrorMsg  = "The page size must be a number greater than 0."

	InvalidSortCriteriaErrorCode = "INVALID_SORT_CRITERIA"
	InvalidSortCriteriaErrorMsg  = "The sort criteria is invalid."
)
