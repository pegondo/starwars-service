package client

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pegondo/starwars/service/internal/errors"
	"github.com/pegondo/starwars/service/internal/handler"
	"github.com/pegondo/starwars/service/internal/resources/swapi"
)

// RetrievePlanets calls the retrieve planets endpoint from c.addr and returns its
// response. If the response is valid, its data will be in planetsResp. If the
// response is invalid and the response has errors.ResponseError format, err
// will contain that error; otherwise, err will contain the returned error.
func (c *Client) RetrievePlanets(opts requestOpts) (planetsResp Response[swapi.Planet], err error) {
	url := c.buildUrl("planets", opts)

	resp, err := http.Get(url)
	if err != nil {
		return planetsResp, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return planetsResp, err
	}

	if isStatusOk(resp.StatusCode) {
		// If the response is a success response, return it.
		var apiResp handler.Response[swapi.Planet]
		if err = json.Unmarshal(bodyBytes, &apiResp); err == nil {
			return Response[swapi.Planet]{
				StatusCode: resp.StatusCode,
				Response:   apiResp,
			}, nil
		}
	} else {
		// If the response is an error response, return it as an error.
		var errResp *errors.ResponseError
		if err = json.Unmarshal(bodyBytes, &errResp); err == nil {
			return planetsResp, errResp
		}
	}

	return planetsResp, ErrInvalidResponseFormat
}
