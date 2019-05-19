import { gql } from 'apollo-boost';
import { Query } from 'react-apollo';

import { Repo, RepoVariables } from './types/Repo';

export const REPO_QUERY = gql`
query Repo($owner: String!, $name: String!, $host: RepositoryHost) {
  repo(host: $host, owner: $owner, name: $name) {
    repository {
      id
      description
    }
    burndown(type: GLOBAL) {
      ... on GlobalBurndown {
        type
        entries {
          start
          bands
        }
      }
      ... on FileBurndown {
        type
      }
      ... on AuthorBurndown {
        type
      }
      ... on BurndownAlert {
        type
        alert
      }
    }
  }
}`;

export class RepoQuery extends Query<Repo, RepoVariables> { }
