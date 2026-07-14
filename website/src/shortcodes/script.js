const { assetHash } = require('../_data/assetHash.js');

module.exports = ({ src, defer = true, isModule = true, ...props }) => {
  return `
    <script
      src="/scripts/${src}?${assetHash}"
      ${isModule ? 'type="module"' : ''}
      ${defer ? 'defer' : ''}
      ${Object.keys(props)
        .filter((key) => key !== '__keywords')
        .map((key) => `${key}="${props[key]}"`)
        .join(' ')}
    ></script>
  `;
};
