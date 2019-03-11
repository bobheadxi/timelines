CREATE TABLE repositories (
  id              SERIAL PRIMARY KEY,
  installation_id TEXT UNIQUE,
  owner           TEXT,
  name            TEXT,
  UNIQUE(owner, name)
);

CREATE TYPE github_item_type AS ENUM ('unknown', 'issue', 'pull_request');
CREATE TABLE github_items (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  github_id    INTEGER,
  number       INTEGER,
  type         github_item_type,

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
  date         DATE,
  PRIMARY KEY (fk_repo_id, date),

  global_delta INTEGER
);

CREATE TABLE git_burndowns_files (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  filename     TEXT,
  date         DATE,
  PRIMARY KEY (fk_repo_id, filename, date),

  delta        INTEGER
);

CREATE TABLE git_burndowns_contributors (
  fk_repo_id   INTEGER REFERENCES repositories(id) ON DELETE CASCADE,
  contributor  TEXT,
  date         DATE,
  PRIMARY KEY (fk_repo_id, contributor, date),

  delta        INTEGER
);
