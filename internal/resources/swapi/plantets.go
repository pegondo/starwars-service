package swapi

import (
	"time"

	internalRequest "github.com/pegondo/starwars-service/internal/request"
)

// planetsEndpoint is the endpoint to request for planets in SWAPI.
const planetsEndpoint = "planets"

// Planet represents a planet in SWAPI.
// Source: https://swapi.dev/documentation#planets
type Planet struct {
	// Name is the name of the planet.
	Name string `json:"name"`
	// Diameter is the diameter of the planet in kilometers.
	Diameter string `json:"diameter"`
	// RotationPeriod is the number of hours it takes for the planet to complete
	// a single orbit of its axis.
	RotationPeriod string `json:"rotation_period"`
	// OrbitalPeriod is the number of days it takes for the planet to complete a
	// single orbit of its local star.
	OrbitalPeriod string `json:"orbital_period"`
	// Gravity is a number denoting the gravity of this planet, where "1" is
	// normal or 1 standard G. "2" is twice or 2 standard Gs. "0.5" is half or
	// 0.5 standard Gs.
	Gravity string `json:"gravity"`
	// Population is the average population of sentient beings inhabiting the
	// planet.
	Population string `json:"population"`
	// Climate is the climate of this planet. Comma separated if diverse.
	Climate string `json:"climate"`
	// Terrain is the terrain of this planet. Comma separated if diverse.
	Terrain string `json:"terrain"`
	// SurfaceWater is the percentage of the planet surface that is naturally
	// occurring water or bodies of water.
	SurfaceWater string `json:"surface_water"`
	// Url is the URL to the resource of this planet.
	Url string `json:"url"`
	// Created is the time when the resource of this planet was created.
	Created time.Time `json:"created"`
	// Edited is the time when the resource of this planet was edited for the
	// last time.
	Edited time.Time `json:"edited"`
}

// GetName returns the planet name.
func (p Planet) GetName() string {
	return p.Name
}

// GetCreated returns the planet's resouce creation time in SWAPI.
func (p Planet) GetCreated() time.Time {
	return p.Created
}

// RetrievePlanets requests the SWAPI for planets. The SWAPI doesn't support
// pagination with variable page sizes, but this function does the maths and
// requests the endpoint various times if needed to return the data for the
// given page and page size. If params.Search is not "", the planets returned
// will contain the value of search in their name. If params.SortCriteria isn't
// nil, the planets will be ordered with the defined criteria.
func RetrievePlanets(
	params internalRequest.RequestParams,
) (
	planetsResp SwapiResponse[Planet],
	err error,
) {
	if params.SortCriteria != nil {
		return retrieveAllAndSort[Planet](planetsEndpoint, params)
	}
	return retrievePage[Planet](planetsEndpoint, params)
}
