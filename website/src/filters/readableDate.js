const { DateTime } = require("luxon");
require("dotenv").config();

module.exports = (value, format = { month: "long", day: "numeric" }) => {
  const dateObject = DateTime.fromISO(value);
  return dateObject.setLocale(process.env.LOCALE).toLocaleString(format);
};
