import React, { Component, ReactElement } from 'react';
import { match } from 'react-router-dom';

import Loading from 'components/Loading/Loading';
import { Error } from 'components/alerts';
import { Burndown } from 'components/vis';

import { getHostTypeFromHost } from 'lib';
import { RepoQuery, REPO_QUERY } from 'lib/queries/repo';
import { BurndownType } from 'lib/queries/types/global';

interface TimelineQuery {
  host: string;
  owner: string;
  name: string;
}

// https://getuikit.com/docs/overlay

class Timeline extends Component<{
  match: match<TimelineQuery>;
}> {
  public render(): ReactElement {
    const { match: { params: { host, owner, name } } } = this.props;

    return (
      <div>
        <RepoQuery
          query={REPO_QUERY}
          variables={{
            owner, name, host: getHostTypeFromHost(host), type: BurndownType.FILE,
          }}
        >
          {({ loading, error, data }): ReactElement => {
            // deal with random edge cases
            if (loading) return <Loading />;
            if (error) return <Error message={`Error :( ${error.message}`} />;
            if (!data || !data.repo) return <Error message="no data found" />;
            const { repo: { repository, burndown } } = data;
            if (!burndown) return <Error message="no data found" />;

            // render visualisation
            return (
              <div>
                <div className="margin-sides-48">
                  <h1 className="uk-heading-line uk-text-center margin-sides-l">
                    <span>{`${host}/${owner}/${name}`}</span>
                  </h1>
                </div>
                <h3>{repository.description}</h3>
                <Burndown data={burndown} />
              </div>
            );
          }}
        </RepoQuery>
      </div>
    );
  }
}

export default Timeline;
