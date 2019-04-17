// Configuration for the Apollo command line tool.
// Reference: https://www.apollographql.com/docs/references/apollo-config
module.exports = {
  includes: ['src/lib/queries/*.tsx'],
  client: {
    name: 'Timelines',
    service: {
      name: 'timelines-api',
      localSchemaFile: '../graphql/schema.graphql',
    },
  },
};
