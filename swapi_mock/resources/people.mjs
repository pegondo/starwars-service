import {
  getPseudoRandomString,
  getPseudoRandomElement,
  getPseudoRandomDate,
} from "./pseudoRandomGenerator.mjs";

const NUM_ELEMENTS = 55;

const VALID_COLORS = ["brown", "blue", "green", "black"];

const VALID_GENDERS = ["male", "female", "unknown"];

const VALID_HEIGHTS = [
  "160",
  "165",
  "170",
  "175",
  "180",
  "185",
  "190",
  "195",
  "200",
];

const VALID_MASSES = ["60", "65", "70", "75", "80", "85", "90", "95", "100"];

const PEOPLE = [...new Array(NUM_ELEMENTS)].map((_, index) => {
  const name = getPseudoRandomString("Name", index);
  const birth_year = getPseudoRandomDate(index);
  const eye_color = getPseudoRandomElement(VALID_COLORS, index);
  const gender = getPseudoRandomElement(VALID_GENDERS, index);
  const hair_color = getPseudoRandomElement(VALID_COLORS, index);
  const height = getPseudoRandomElement(VALID_HEIGHTS, index);
  const mass = getPseudoRandomElement(VALID_MASSES, index);
  const skin_color = getPseudoRandomElement(VALID_COLORS, index);
  const url = getPseudoRandomString("people", index);
  const created = getPseudoRandomDate(index);
  const edited = getPseudoRandomDate(index);
  return {
    name,
    birth_year,
    eye_color,
    gender,
    hair_color,
    height,
    mass,
    skin_color,
    url,
    created,
    edited,
  };
});

export default PEOPLE;
