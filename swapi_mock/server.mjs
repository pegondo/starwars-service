import express from "express";
import PEOPLE from "./resources/people.mjs";
import PLANETS from "./resources/planets.mjs";

const BASE_URL = "http://localhost:3000";

// SWAPI had a fixed page size.
const PAGE_SIZE = 10;

const app = express();
const port = 3000;

const getPageQueryParam = (req) => {
  const pageParam = req.query.page;
  return parseInt(pageParam) || 1;
};

const getSearchQueryParam = (req) => req.query.search;

const getQueryParams = (req) => ({
  page: getPageQueryParam(req),
  search: getSearchQueryParam(req),
});

const filterResources = (resources, search) =>
  search
    ? resources.filter(({ name }) =>
        name.toLowerCase().includes(search.toLowerCase())
      )
    : resources;

const applyPagination = (resources, page) => {
  const res = Object.assign([], resources);

  const startIndex = (page - 1) * PAGE_SIZE;
  const endIndex = Math.min(startIndex + PAGE_SIZE, res.length);

  if (startIndex > res.length) return [];
  return res.slice(startIndex, endIndex);
};

const buildNextUrl = (endpoint, count, page, search) => {
  const numAlreadyDeliveredResources = page * PAGE_SIZE;
  const thereAreResourcesLeft = count > numAlreadyDeliveredResources;
  if (!thereAreResourcesLeft) return undefined;

  let nextUrl = `${BASE_URL}/${endpoint}?page=${page + 1}`;
  if (search !== undefined) {
    nextUrl = `${nextUrl}&search=${search}`;
  }
  return nextUrl;
};

const getCount = (numElements, responseLength) => {
  if (responseLength === 0) {
    // When there are no elements in the response, return 0. This isn't the
    // normal behavior in a REST API, but this is how SWAPI managed this
    // scenario.
    return 0;
  }
  return numElements;
};

const applyQueryParams = (endpoint, resources, { page, search }) => {
  const filteredResources = filterResources(resources, search);

  const paginatedResources = applyPagination(filteredResources, page);
  const count = getCount(filteredResources.length, paginatedResources.length);
  const next = buildNextUrl(endpoint, count, page, search);

  return {
    count,
    next,
    results: paginatedResources,
  };
};

app.get("/people", (req, res) => {
  const params = getQueryParams(req);

  const response = applyQueryParams("people", PEOPLE, params);

  res.json(response);
});

app.get("/planets", (req, res) => {
  const params = getQueryParams(req);

  const response = applyQueryParams("planets", PLANETS, params);

  res.json(response);
});

app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
