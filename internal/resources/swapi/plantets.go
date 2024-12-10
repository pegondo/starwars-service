package swapi

import (
	"fmt"
	"math"
	"sort"
	"starwars/service/internal/utils"
	"time"
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

// retrieveAllPlanetsAndSort retrieves all the planets in SWAPI and sorts them
// using the given criteria to return the information paginated with the given
// page number and size. If search isn't "", the names of the planets in
// planetsResp.Result will contain the value of search.
func retrieveAllPlanetsAndSort(
	page,
	pageSize int,
	search string,
	sortCriteria SortCriteria,
) (
	planetsResp SwapiResponse[Planet],
	err error,
) {
	planets, err := retrieveAll[Planet](planetsEndpoint, search)
	if err != nil {
		return planetsResp, err
	}

	var lessFn func(i, j int) bool
	switch sortCriteria.Field {
	case NameSortField:
		lessFn = func(i, j int) bool {
			return planets.Results[i].Name < planets.Results[j].Name
		}

	case CreatedSortField:
		lessFn = func(i, j int) bool {
			return planets.Results[i].Created.Before(planets.Results[j].Created)
		}
	default:
		return planetsResp, fmt.Errorf("invalid sort field '%s'", sortCriteria.Field)
	}
	sort.Slice(planets.Results, lessFn)

	if sortCriteria.Order == DescendingOrder {
		utils.ReverseSlice(planets.Results)
	}

	minIdx := (page - 1) * pageSize
	// TODO: Fix here as well.
	maxIdx := int(math.Min(float64(page*pageSize), float64(len(planets.Results))))
	planets.Results = planets.Results[minIdx:maxIdx]

	return planets, nil
}

// retrievePlanetsPage retrieves the planets from the SWAPI with the given page
// number and size. If search isn't "", all the elements of planetsResp.Results
// will contain the value of search.
func retrievePlanetsPage(
	page,
	pageSize int,
	search string,
) (
	planetsResp SwapiResponse[Planet],
	err error,
) {
	numAlreadyRequestedPlanets := (page - 1) * pageSize
	initialPage := int(numAlreadyRequestedPlanets/swapiPageSize) + 1
	initialPageOffset := numAlreadyRequestedPlanets % swapiPageSize

	return retrievePage(SwapiResponse[Planet]{}, planetsEndpoint, search, pageSize, initialPage, initialPageOffset)
}

// RetrievePlanets requests the SWAPI for planets. The SWAPI doesn't support
// pagination with variable page sizes, but this function does the maths and
// requests the endpoint various times if needed to return the data for the
// given page and page size. If search is not "", the planets returned will
// contain the value of search in their name.
func RetrievePlanets(
	page,
	pageSize int,
	search string,
	sortCriteria *SortCriteria,
) (
	planetsResp SwapiResponse[Planet],
	err error,
) {
	if sortCriteria != nil {
		return retrieveAllPlanetsAndSort(page, pageSize, search, *sortCriteria)
	}
	return retrievePlanetsPage(page, pageSize, search)
}
