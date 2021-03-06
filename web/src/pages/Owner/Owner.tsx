import React, { Component, ReactElement } from 'react';
import { match } from 'react-router-dom';

import Loading from 'components/Loading/Loading';
import CardSet, { Card } from 'components/CardSet/CardSet';
import { Error } from 'components/alerts';

import { getHostTypeFromHost } from 'lib';
import { ReposQuery, REPOS_QUERY } from 'lib/queries/repos';

interface OwnerQuery {
  host: string;
  owner: string;
}

class Owner extends Component<{
  match: match<OwnerQuery>;
}> {
  public render(): ReactElement {
    const { match: { params: { host, owner } } } = this.props;

    return (
      <div>
        <div className="margin-sides-48">
          <h1 className="uk-heading-line uk-text-center pad-bot-48">
            <span>{`${host}/${owner}`}</span>
          </h1>

          <ReposQuery query={REPOS_QUERY} variables={{ owner, host: getHostTypeFromHost(host) }}>
            {({ loading, error, data }): ReactElement => {
              if (loading) return <Loading />;
              if (error) return <Error message={`Error :( ${error.message}`} />;

              if (!data || !data.repos) return <Error message="no data found" />;
              const { repos } = data;

              return (
                <CardSet cards={repos.map((r): Card => ({
                  title: r.name,
                  body: r.description,
                  button: {
                    href: `${host}/${owner}`,
                    text: 'View Project',
                  },
                }))}
                />
              );
            }}
          </ReposQuery>
        </div>
      </div>
    );
  }
}

export default Owner;
