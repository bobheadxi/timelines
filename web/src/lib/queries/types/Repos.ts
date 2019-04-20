/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

import { RepositoryHost } from "./global";

// ====================================================
// GraphQL query operation: Repos
// ====================================================

export interface Repos_repos {
  __typename: "Repository";
  id: number;
  name: string;
  description: string;
}

export interface Repos {
  repos: Repos_repos[] | null;
}

export interface ReposVariables {
  owner: string;
  host?: RepositoryHost | null;
}
