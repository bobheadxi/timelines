/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

import { RepositoryHost, BurndownType } from "./global";

// ====================================================
// GraphQL query operation: Repo
// ====================================================

export interface Repo_repo_repository {
  __typename: "Repository";
  id: number;
  description: string;
}

export interface Repo_repo_burndown_GlobalBurndown_entries {
  __typename: "BurndownEntry";
  start: any;
  bands: number[];
}

export interface Repo_repo_burndown_GlobalBurndown {
  __typename: "GlobalBurndown";
  type: BurndownType;
  entries: Repo_repo_burndown_GlobalBurndown_entries[] | null;
}

export interface Repo_repo_burndown_FileBurndown_file_entry {
  __typename: "BurndownEntry";
  start: any;
  bands: number[];
}

export interface Repo_repo_burndown_FileBurndown_file {
  __typename: "FileBurndownEntry";
  file: string;
  entry: Repo_repo_burndown_FileBurndown_file_entry;
}

export interface Repo_repo_burndown_FileBurndown {
  __typename: "FileBurndown";
  type: BurndownType;
  file: Repo_repo_burndown_FileBurndown_file[] | null;
}

export interface Repo_repo_burndown_AuthorBurndown {
  __typename: "AuthorBurndown";
  type: BurndownType;
}

export interface Repo_repo_burndown_BurndownAlert {
  __typename: "BurndownAlert";
  type: BurndownType;
  alert: string;
}

export type Repo_repo_burndown = Repo_repo_burndown_GlobalBurndown | Repo_repo_burndown_FileBurndown | Repo_repo_burndown_AuthorBurndown | Repo_repo_burndown_BurndownAlert;

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
  type?: BurndownType | null;
}
