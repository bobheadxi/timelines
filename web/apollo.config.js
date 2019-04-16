module.exports = {
  includes: [ 'src/queries/*.tsx' ],
  client: {
    name: 'Timelines',
    service: {
      name: 'timelines-api',
      localSchemaFile: '../graphql/schema.graphql',
    }
  }
};
