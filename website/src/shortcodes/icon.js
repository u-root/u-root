module.exports = (icon) => {
  let style = '';
  let srText = "aria-hidden='true' focusable='false'";

  if (icon.width) {
    style += 'width:' + icon.width + '; height:' + icon.width + ';';
  }
  if (icon.alt) {
    srText = "aria-label='" + icon.alt + "'";
  }
  return `<svg style="${style}" class="icon" ${srText}><use xlink:href="/icons/icons.svg#svg-${icon.icon}"/></svg>`;
};
