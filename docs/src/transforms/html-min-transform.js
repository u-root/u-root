const htmlmin = require("html-minifier");

module.exports = (value, outputPath) => {
  if (outputPath && outputPath.endsWith(".html")) {
    return htmlmin.minify(value, {
      useShortDoctype: true,
      removeComments: true,
      collapseWhitespace: true,
      minifyCSS: true,
      minifyJs: true,
    });
  }

  return value;
};
