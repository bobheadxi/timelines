import { gql } from 'apollo-boost'

export const REPOS_QUERY = gql`
query GetRepos($owner: String!) {
  repos(owner: $owner) {
    id
    name
  }
}`;
