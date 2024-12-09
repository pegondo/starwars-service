package swapi

import (
	"fmt"
	"math"
	"time"
)

// peopleResource is the endpoint to request for people in SWAPI.
const peopleResource = "people"

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

// TODO: Remove the unneeded fields.

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
	// Films is an array of film resource URLs that this person has been in.
	Films []string `json:"films"`
	// Species is an array of species resource URLs that this person belongs to.
	Species []string `json:"species"`
	// Starships is an array of starship resource URLs that this person has
	// piloted.
	Starships []string `json:"starships"`
	// Vehicles is an array of vehicle resource URLs that this person has
	// piloted.
	Vehicles []string `json:"vehicles"`
	// Url is the URL to the resource of this person.
	Url string `json:"url"`
	// Created is the time when the resource of this person was created.
	Created time.Time `json:"created"`
	// Edited is the time when the resource of this person was edited for the
	// last time.
	Edited time.Time `json:"edited"`
}

// PeopleResponse represents the response of SWAPI to a retrieve all people
// request.
type PersonResponse struct {
	// Count represents the number of elements in the people collection.
	Count int `json:"count"`
	// Results is the list of people that came in the response.
	Results []Person `json:"results"`
}

// retrievePeople is a recursive solution to request the people SWAPI endpoint
// given a variable page size.
func retrievePeople(people PersonResponse, remainingPeople, pageNumber, offset int) (peopleResp PersonResponse, err error) {
	if remainingPeople <= 0 {
		return people, nil
	}

	peopleResp, err = request(fmt.Sprintf("%s/%s?page=%d", swapiBaseUrl, peopleResource, pageNumber))
	if err != nil {
		return peopleResp, fmt.Errorf("error while requesting the people endpoint :: %v", err)
	}
	if peopleResp.Count == 0 {
		// If there are no results, return an empty person response.
		return PersonResponse{
			Count:   0,
			Results: []Person{},
		}, nil
	}

	remainingPeople = int(math.Min(float64(remainingPeople), float64(peopleResp.Count)))

	minIdx := offset
	maxIdx := int(math.Min(swapiPageSize, float64(minIdx+remainingPeople)))
	peopleResp.Results = append(people.Results, peopleResp.Results[minIdx:maxIdx]...)
	numElementsAdded := maxIdx - minIdx
	return retrievePeople(peopleResp, remainingPeople-numElementsAdded, pageNumber+1, 0)
}

// RetrievePeople requests the SWAPI for people. The SWAPI doesn't support
// pagination with variable page sizes, but this function does the maths and
// requests the endpoint various times if needed to return the data for the
// given page and page size.
func RetrievePeople(page, pageSize int) (peopleResp PersonResponse, err error) {
	numAlreadyRequestedPeople := (page - 1) * pageSize
	initialPage := int(numAlreadyRequestedPeople/swapiPageSize) + 1
	initialPageOffset := numAlreadyRequestedPeople % swapiPageSize

	return retrievePeople(PersonResponse{}, pageSize, initialPage, initialPageOffset)
}
