/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

import { RepositoryHost } from "./global";

// ====================================================
// GraphQL query operation: Repo
// ====================================================

export interface Repo_repo_repository {
  __typename: "Repository";
  id: number;
  description: string;
}

export interface Repo_repo_burndown_AuthorBurndown {
  __typename: "AuthorBurndown" | "FileBurndown";
}

export interface Repo_repo_burndown_GlobalBurndown_entries {
  __typename: "BurndownEntry";
  start: any;
  bands: number[];
}

export interface Repo_repo_burndown_GlobalBurndown {
  __typename: "GlobalBurndown";
  entries: Repo_repo_burndown_GlobalBurndown_entries[] | null;
}

export interface Repo_repo_burndown_BurndownAlert {
  __typename: "BurndownAlert";
  alert: string;
}

export type Repo_repo_burndown = Repo_repo_burndown_AuthorBurndown | Repo_repo_burndown_GlobalBurndown | Repo_repo_burndown_BurndownAlert;

export interface Repo_repo {
  __typename: "RepositoryAnalytics";
  repository: Repo_repo_repository;
  burndown: Repo_repo_burndown | null;
}

export interface Repo {
  repo: Repo_repo | null;
}

export interface RepoVariables {
  owner: string;
  name: string;
  host?: RepositoryHost | null;
}
