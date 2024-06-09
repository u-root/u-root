const markdownIt = require('markdown-it')({
  html: true,
  breaks: true,
  linkify: true,
}).use(require('markdown-it-footnote'));

// Filters
const readableDate = require('./src/filters/readableDate.js');
const w3DateFilter = require('./src/filters/w3-date-filter.js');
const markdownFilter = require('./src/filters/markdown-filter.js');

// Plugins
const svgSprite = require('eleventy-plugin-svg-sprite');
const { globPlugin } = require('esbuild-plugin-glob');

// Transforms
const htmlMinTransform = require('./src/transforms/html-min-transform.js');
const purgeCSS = require('./src/transforms/css-purge-inline.js');

const eleventyBuildSystem = require('@cagov/11ty-build-system');

// Create a helpful production flag
const isProduction = process.env.NODE_ENV === 'production';

const path = require('path');

module.exports = (eleventyConfig) => {
  // Set directories to pass through to the dist folder
  eleventyConfig.addPassthroughCopy('./src/images/');
  eleventyConfig.addPassthroughCopy('./src/fonts/');

  // Add filters
  eleventyConfig.addFilter('readableDate', readableDate);
  eleventyConfig.addFilter('w3DateFilter', w3DateFilter);
  eleventyConfig.addFilter('limit', function (arr, limit) {
    return arr.slice(0, limit);
  });
  eleventyConfig.addFilter('markdownFilter', markdownFilter);

  // Configure markdown options
  markdownIt.renderer.rules.footnote_block_open = () =>
    '<footer>\n' +
    '<h4 class="sr-only">Footnotes</h4>\n' +
    '<ol class="footnotes-list">\n';

  eleventyConfig.setLibrary('md', markdownIt);

  // Add Shortcodes
  eleventyConfig.addShortcode('icon', require('./src/shortcodes/icon.js'));
  eleventyConfig.addShortcode('script', require('./src/shortcodes/script.js'));

  // Plugins
  eleventyConfig.addPlugin(svgSprite, {
    path: './src/icons', // relative path to SVG directory
    outputFilepath: './dist/icons/icons.svg',
  });

  eleventyConfig.addPlugin(eleventyBuildSystem, {
    processors: {
      esbuild: {
        watch: ['src/scripts/**/*'],
        options: {
          entryPoints: [path.resolve(__dirname, 'src/scripts/**/*')],
          bundle: true,
          minify: isProduction,
          outdir: 'dist/scripts',
          splitting: true,
          format: 'esm',
          plugins: [globPlugin()],
        },
      },
    },
  });

  // Only minify HTML if we are in production because it slows builds _right_ down
  if (isProduction) {
    eleventyConfig.addTransform('htmlmin', htmlMinTransform);
    eleventyConfig.addTransform('purgeCSS', purgeCSS);
  }

  // Tell 11ty to use the .eleventyignore and ignore our .gitignore file
  eleventyConfig.setUseGitIgnore(false);

  return {
    markdownTemplateEngine: 'njk',
    dataTemplateEngine: 'njk',
    htmlTemplateEngine: 'njk',
    dir: {
      input: 'src',
      output: 'dist',
    },
  };
};
