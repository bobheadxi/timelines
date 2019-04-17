import { RepositoryHost } from './queries/types/global';

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
};
