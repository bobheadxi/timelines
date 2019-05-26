import { RepositoryHost, BurndownType } from './queries/types/global';

// Alias some generated types for readability
/* eslint-disable @typescript-eslint/camelcase */
import {
  Repo_repo_burndown,
  Repo_repo_burndown_GlobalBurndown,
  Repo_repo_burndown_GlobalBurndown_entries,

  Repo_repo_burndown_FileBurndown,
  Repo_repo_burndown_FileBurndown_file,
  Repo_repo_burndown_FileBurndown_file_entry,

  Repo_repo_burndown_AuthorBurndown,

  Repo_repo_burndown_BurndownAlert,
} from './queries/types/Repo';

export type RepoBurndown = Repo_repo_burndown;

export type GlobalBurndown = Repo_repo_burndown_GlobalBurndown;
export type GlobalBurndownEntry = Repo_repo_burndown_GlobalBurndown_entries;

export type FilesBurndown = Repo_repo_burndown_FileBurndown;
export type FilesBurndownFile = Repo_repo_burndown_FileBurndown_file;
export type FilesBurndownEntry = Repo_repo_burndown_FileBurndown_file_entry;

export type AuthorBurndown = Repo_repo_burndown_AuthorBurndown;

export type BurndownAlert = Repo_repo_burndown_BurndownAlert;
/* eslint-enable @typescript-eslint/camelcase */

function getHostTypeFromHost(host: string): RepositoryHost {
  switch (host) {
    case 'github.com': return RepositoryHost.GITHUB;
    case 'gitlab.com': return RepositoryHost.GITLAB;
    case 'bitbucket.com': return RepositoryHost.BITBUCKET;
    default: throw new Error(`invalid code host '${host}'`);
  }
}

export {
  getHostTypeFromHost,
  BurndownType,
};
