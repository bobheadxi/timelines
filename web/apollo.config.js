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
