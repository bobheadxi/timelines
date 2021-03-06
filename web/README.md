# Timelines Web [![Netlify Status](https://api.netlify.com/api/v1/badges/b56788d9-0743-4b39-a307-66e2c99bd428/deploy-status)](https://app.netlify.com/sites/timelines-bobheadxi/deploys)

This document outlines development notes for Timelines' web interface. It is
currently deployed using [Netlify](https://www.netlify.com/).

* [Code Style](#code-style)
* [GraphQL Queries](#graphql-queries)
* [Styling](#styling)
  * [UIkit](#uikit)
  * [SCSS](#scss)

## Code Style

Styles are enforced by [eslint](https://eslint.org/), with rules defined in
[`./.eslintrc.js`](./.eslintrc.js). 

Note that because [`create-react-app`](https://github.com/facebook/create-react-app)
is silly, `eslint` must be installed *globally* (and not in the project's
`node_modules`) before running the linter:

```
npm i -g eslint
npm run lint
```

## GraphQL Queries

The web app connects to the backend server using GraphQL, with a client powered by
[Apollo](https://github.com/apollographql/apollo-client). Relevant documentation:
[link](https://www.apollographql.com/docs/react/).

The API schema is defined in [`../graphql/schema.graphql`](../graphql/schema.graphql).

GraphQL API Typescript definitions for the web app are generated using
[Apollo-Tooling](https://github.com/apollographql/apollo-tooling), though the
generated code still needs a small wrapper layer on top - see
[`./src/lib/queries/repos.tsx`](src/lib/queries/repos.tsx). To update the generated
code, run:

```
npm run graphql
```

Settings for code generation are defined in [`./apollo.config.js`](./apollo.config.js).

## Styling

### UIkit

The Timelines web app leverages [uikit](https://getuikit.com/docs/introduction)
for most UI elements.

Certain UIkit styles, such as [`uk-grid`](https://getuikit.com/docs/grid),
require that you add a custom attribute to an HTML element.
[React unfortunately ignores non-standard attributes](https://zhenyong.github.io/react/docs/jsx-gotchas.html#custom-html-attributes),
so such attributes must be prefixed with `data-`. For example:

```jsx
<div className="uk-child-width-1-2@s uk-grid-match" data-uk-grid></div>
```

### SCSS

Custom style classes are defined in [`./src/styles`](./src/styles/_all.scss).
