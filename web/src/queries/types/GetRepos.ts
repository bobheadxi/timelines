/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: GetRepos
// ====================================================

export interface GetRepos_repos {
  __typename: "Repository";
  id: number;
  name: string;
}

export interface GetRepos {
  repos: GetRepos_repos[] | null;
}

export interface GetReposVariables {
  owner: string;
}
