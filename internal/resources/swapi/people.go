package swapi

import (
	"fmt"
	"math"
	"sort"
	"starwars/service/internal/utils"
	"time"
)

// peopleEndpoint is the endpoint to request for people in SWAPI.
const peopleEndpoint = "people"

// Gender represents a person gender.
type Gender string

const (
	// MaleGender represents the male gender.
	MaleGender = "Male"
	// FemaleGender represents the female gender.
	FemaleGender = "Female"
	// UnknownGender represents an unknown gender.
	UnknownGender = "Unknown"
)

// Person is the data structure SWAPI uses to define a person.
// Source: https://swapi.dev/documentation#people
type Person struct {
	// The name of the person.
	Name string `json:"name"`
	// BirthYear is the birth year of the person, using the in-universe standard
	// of BBY or ABY - Before the Battle of Yavin or After the Battle of Yavin.
	// The Battle of Yavin is a battle that occurs at the end of Star Wars
	// episode IV: A New Hope.
	BirthYear string `json:"birth_year"`
	// EyeColor is the eye color of this person. Will be "unknown" if not known
	// or "n/a" if the person does not have an eye.
	EyeColor *string `json:"eye_color"`
	// Gender is the gender of the person.
	Gender *Gender `json:"gender"`
	// HairColor is the color of the person's hair.
	HairColor *string `json:"hair_color"`
	// Height is the height of the person in centimeters.
	Height string `json:"height"`
	// Mass is the mass of the person in kilograms.
	Mass string `json:"mass"`
	// SkinColor is the color of the person.
	SkinColor string `json:"skin_color"`
	// Homeworld is the URL of a planet resource, a planet that this person was
	// born on or inhabits.
	Homeword string `json:"homeworld"`
	// Url is the URL to the resource of this person.
	Url string `json:"url"`
	// Created is the time when the resource of this person was created.
	Created time.Time `json:"created"`
	// Edited is the time when the resource of this person was edited for the
	// last time.
	Edited time.Time `json:"edited"`
}

// retrieveAllPeopleAndSort retrieves all the people in SWAPI and sorts them
// using the given criteria to return the information paginated with the given
// page number and size. If search isn't "", the names of the people in
// peopleResp.Result will contain the value of search.
func retrieveAllPeopleAndSort(
	page,
	pageSize int,
	search string,
	sortCriteria SortCriteria,
) (
	peopleResp SwapiResponse[Person],
	err error,
) {
	people, err := retrieveAll[Person](peopleEndpoint, search)
	if err != nil {
		return peopleResp, err
	}

	var lessFn func(i, j int) bool
	switch sortCriteria.Field {
	case NameSortField:
		lessFn = func(i, j int) bool {
			return people.Results[i].Name < people.Results[j].Name
		}

	case CreatedSortField:
		lessFn = func(i, j int) bool {
			return people.Results[i].Created.Before(people.Results[j].Created)
		}
	default:
		return peopleResp, fmt.Errorf("invalid sort field '%s'", sortCriteria.Field)
	}
	sort.Slice(people.Results, lessFn)

	if sortCriteria.Order == DescendingOrder {
		utils.ReverseSlice(people.Results)
	}

	minIdx := (page - 1) * pageSize
	if minIdx > len(people.Results) {
		return SwapiResponse[Person]{
			Count:   people.Count,
			Results: []Person{},
		}, nil
	}
	maxIdx := int(math.Min(float64(page*pageSize), float64(len(people.Results))))
	people.Results = people.Results[minIdx:maxIdx]

	return people, nil
}

// retrievePeoplePage retrieves the people from the SWAPI with the given page
// number and size. If search isn't "", all the elements of peopleResp.Results
// will contain the value of search.
func retrievePeoplePage(
	page,
	pageSize int,
	search string,
) (
	peopleResp SwapiResponse[Person],
	err error,
) {
	numAlreadyRequestedPeople := (page - 1) * pageSize
	initialPage := int(numAlreadyRequestedPeople/swapiPageSize) + 1
	initialPageOffset := numAlreadyRequestedPeople % swapiPageSize

	return retrievePage(SwapiResponse[Person]{}, peopleEndpoint, search, pageSize, initialPage, initialPageOffset)
}

// RetrievePeople requests the SWAPI for people. The SWAPI doesn't support
// pagination with variable page sizes, but this function does the maths and
// requests the endpoint various times if needed to return the data for the
// given page and page size. If search is not "", the people returned will
// contain the value of search in their name.
func RetrievePeople(
	page,
	pageSize int,
	search string,
	sortCriteria *SortCriteria,
) (
	peopleResp SwapiResponse[Person],
	err error,
) {
	if sortCriteria != nil {
		return retrieveAllPeopleAndSort(page, pageSize, search, *sortCriteria)
	}
	return retrievePeoplePage(page, pageSize, search)
}
