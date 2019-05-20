import { IntrospectionFragmentMatcher } from 'apollo-cache-inmemory';
import fragmentData from './types/fragmentTypes.json';

export const fragmentMatcher = new IntrospectionFragmentMatcher({
  introspectionQueryResultData: fragmentData,
});
