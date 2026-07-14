module.exports = (value) => {
  const dateObject = new Date(value);

  return dateObject.toISOString();
};
