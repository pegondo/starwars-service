package swapi

import (
	"time"

	internalRequest "github.com/pegondo/starwars-service/internal/request"
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
	// Url is the URL to the resource of this person.
	Url string `json:"url"`
	// Created is the time when the resource of this person was created.
	Created time.Time `json:"created"`
	// Edited is the time when the resource of this person was edited for the
	// last time.
	Edited time.Time `json:"edited"`
}

// GetName returns the person name.
func (p Person) GetName() string {
	return p.Name
}

// GetCreated returns the person's resouce creation time in SWAPI.
func (p Person) GetCreated() time.Time {
	return p.Created
}

// RetrievePeople requests the SWAPI for people. The SWAPI doesn't support
// pagination with variable page sizes, but this function does the maths and
// requests the endpoint various times if needed to return the data for the
// given page and page size. If params.Search is not "", the people returned
// will contain the value of search in their name. If params.SortCriteria isn't
// nil, the people will be ordered with the defined criteria.
func RetrievePeople(
	params internalRequest.RequestParams,
) (
	peopleResp SwapiResponse[Person],
	err error,
) {
	if params.SortCriteria != nil {
		return retrieveAllAndSort[Person](peopleEndpoint, params)
	}
	return retrievePage[Person](peopleEndpoint, params)
}
