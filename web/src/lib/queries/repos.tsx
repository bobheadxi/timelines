import { gql } from 'apollo-boost';
import { Query } from 'react-apollo';

import { Repos, ReposVariables } from './types/Repos';

export const REPOS_QUERY = gql`
query Repos($owner: String!, $host: RepositoryHost) {
  repos(owner: $owner, host: $host) {
    id
    name
  }
}`;

export class ReposQuery extends Query<Repos, ReposVariables> { }
