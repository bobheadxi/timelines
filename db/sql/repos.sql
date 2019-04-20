CREATE TYPE host_service AS ENUM ('unknown', 'github', 'gitlab');
CREATE TABLE repositories (
  id              SERIAL PRIMARY KEY,
  installation_id TEXT UNIQUE,
  host_type       host_service,
  owner           TEXT,
  name            TEXT,
  description     TEXT,
  service_stats   JSONB NULL,

  UNIQUE(owner, name)
);

CREATE TYPE host_item_type AS ENUM ('unknown', 'issue', 'pull_request');
CREATE TABLE host_items (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  host_type    host_service,
  host_id      INTEGER,
  type         host_item_type,
  number       INTEGER,

  author       TEXT,
  open_date    DATE,
  close_date   DATE NULL,

  title        TEXT,
  body         TEXT,
  labels       TEXT[],
  reactions    JSONB NULL,

  details      JSONB NULL
);

CREATE TABLE git_burndowns_globals (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  interval     TSRANGE,
  PRIMARY KEY (fk_repo_id, interval),

  delta_bands  BIGINT[]
);

CREATE TABLE git_burndowns_files (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  filename     TEXT,
  interval     TSRANGE,
  PRIMARY KEY (fk_repo_id, filename, interval),

  delta_bands  BIGINT[]
);

CREATE TABLE git_burndowns_contributors (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  contributor  TEXT,
  interval     TSRANGE,
  PRIMARY KEY (fk_repo_id, contributor, interval),

  delta_bands  BIGINT[]
);
