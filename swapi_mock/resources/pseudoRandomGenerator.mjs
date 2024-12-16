export const getPseudoRandomElement = (validElements, index) => {
  const i = index % validElements.length;
  return validElements[i];
};

export const getPseudoRandomString = (label, index) => `${label}-${index}`;

export const getPseudoRandomDate = (index) =>
  new Date(new Date().setHours(new Date().getHours() + index));
