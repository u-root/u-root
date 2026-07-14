const postcssPresetEnv = require("postcss-preset-env");

module.exports = {
  plugins: [
    postcssPresetEnv({
      stage: 0,
      features: {
        "custom-properties": false
      }
    })
  ]
};
