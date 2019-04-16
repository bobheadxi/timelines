# Timelines Web

[![Netlify Status](https://api.netlify.com/api/v1/badges/b56788d9-0743-4b39-a307-66e2c99bd428/deploy-status)](https://app.netlify.com/sites/timelines-bobheadxi/deploys)

This document outlines development notes for Timelines' web interface. It is
currently deployed using [Netlify](https://www.netlify.com/).

## Queries

The web app connects to the backend server using GraphQL, with a client powered by
[Apollo](https://github.com/apollographql/apollo-client). Relevant documentation:
[link](https://www.apollographql.com/docs/react/).

The API schema is defined in [`../graphql/schema.graphql`](../graphql/schema.graphql).

Typescript definitions for the web app are generated using
[Apollo-Tooling](https://github.com/apollographql/apollo-tooling), though the
generated code still needs a small wrapper layer on top - see
[`src/lib/queries/repos.tsx`](src/lib/queries/repos.tsx). To update the generated
code, run:

```
npm run graphql
```

## Styling

### UIkit

Certain UIkit styles, such as [`uk-grid`](https://getuikit.com/docs/grid),
require that you add a custom attribute to an HTML element.
[React unforuntately ignores non-standard attributes](https://zhenyong.github.io/react/docs/jsx-gotchas.html#custom-html-attributes),
so such attributes must be prefixed with `data-`. For example:

```jsx
<div className="uk-child-width-1-2@s uk-grid-match" data-uk-grid></div>
```
