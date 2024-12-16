import {
  getPseudoRandomString,
  getPseudoRandomElement,
  getPseudoRandomDate,
} from "./pseudoRandomGenerator.mjs";

const NUM_ELEMENTS = 55;

const VALID_DIAMETERS = ["5000", "10000", "15000", "20000", "25000", "30000"];

const VALID_ROTATION_PERIODS = ["21", "22", "23", "24", "25", "26", "27"];

const VALID_ORBITAL_PERIODS = ["350", "355", "360", "365", "370", "375", "380"];

const VALID_GRAVITIES = ["0.5 G", "1 G", "1.5 G", "2 G"];

const VALID_POPULATIONS = ["5B", "7B", "7.5B", "8B", "10B"];

const VALID_CLIMATES = ["dry", "wet", "tropical"];

const VALID_TERRAINS = ["mountain", "beach", "lake", "ocean"];

const VALID_SURFACE_WATER = ["0%", "25%", "50%", "75%", "100%"];

const PLANETS = [...new Array(NUM_ELEMENTS)].map((_, index) => {
  const name = getPseudoRandomString("Name", index);
  const diameter = getPseudoRandomElement(VALID_DIAMETERS, index);
  const rotation_period = getPseudoRandomElement(VALID_ROTATION_PERIODS, index);
  const orbital_period = getPseudoRandomElement(VALID_ORBITAL_PERIODS, index);
  const gravity = getPseudoRandomElement(VALID_GRAVITIES, index);
  const population = getPseudoRandomElement(VALID_POPULATIONS, index);
  const climate = getPseudoRandomElement(VALID_CLIMATES, index);
  const terrain = getPseudoRandomElement(VALID_TERRAINS, index);
  const surface_water = getPseudoRandomElement(VALID_SURFACE_WATER, index);
  const url = getPseudoRandomString("planets", index);
  const created = getPseudoRandomDate(index);
  const edited = getPseudoRandomDate(index);
  return {
    name,
    diameter,
    rotation_period,
    orbital_period,
    gravity,
    population,
    climate,
    terrain,
    surface_water,
    url,
    created,
    edited,
  };
});

export default PLANETS;
