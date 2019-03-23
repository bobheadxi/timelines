# Timelines Web

This document outlines development notes for Timelines' web interface.

## CSS

### UIkit

Certain UIkit styles, such as [`uk-grid`](https://getuikit.com/docs/grid),
require that you add a custom attribute to an HTML element.
[React unforuntately ignores non-standard attributes](https://zhenyong.github.io/react/docs/jsx-gotchas.html#custom-html-attributes),
so such attributes must be prefixed with `data-`. For example:

```jsx
<div className="uk-child-width-1-2@s uk-grid-match" data-uk-grid></div>
```
