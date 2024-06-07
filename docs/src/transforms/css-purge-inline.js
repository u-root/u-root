const { PurgeCSS } = require('purgecss');
const { JSDOM } = require('jsdom');
const CleanCSS = require('clean-css');

//function to insert css into the DOM
const insertCss = (html, css) => {
  const dom = new JSDOM(html);
  const { document } = dom.window;

  let head = document.getElementsByTagName('head')[0];
  let style = document.createElement('style');
  style.innerHTML = css;
  head.appendChild(style);

  return dom.serialize();
};

module.exports = async (content, outputPath) => {
  if (outputPath && outputPath.endsWith('.html')) {
    //array of css files to combine
    const cssFiles = ['./src/_includes/css/critical.css'];

    // cleanCSSOptions for minification and inlining css, will fix duplicate media queries
    const cleanCSSOptions = {
      level: 1,
    };

    const purgecssResult = await new PurgeCSS().purge({
      content: [
        {
          raw: content,
          extension: 'html',
        },
      ],
      css: cssFiles,
      safelist: {
        standard: [/^md/, /^lg/, /^code/, /^comment/],
        greedy: [/role$/, /where$/, /is$/, /youtube-embed$/, /vimeo-embed$/],
      },
      keyframes: true,
    });

    let cssMerge = '';

    if (purgecssResult.length > 0) {
      for (let i = 0; i < purgecssResult.length; i++) {
        cssMerge = cssMerge.concat(purgecssResult[i].css);
      }
      const cssMin = new CleanCSS(cleanCSSOptions).minify(cssMerge).styles;

      return insertCss(content, cssMin);
    }
  }
  return content;
};
