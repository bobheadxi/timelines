
function getHostTypeFromHost(host: string): string {
  switch (host) {
  case "github.com": return "GITHUB";
  case "gitlab.com": return "GITLAB";
  case "bitbucket.com": return "BITBUCKET";
  default: throw new Error(`invalid code host '${host}'`)
  }
}

export {
  getHostTypeFromHost,
}
