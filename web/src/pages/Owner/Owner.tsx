import React, { Component, ReactElement } from 'react';
import { match } from 'react-router-dom';

import Loading from '../../components/Loading/Loading';
import CardSet, { Card } from '../../components/CardSet/CardSet';
import { getHostTypeFromHost } from '../../lib';
import { ReposQuery, REPOS_QUERY } from '../../lib/queries/repos';

interface OwnerQuery {
  host: string;
  owner: string;
}

class Owner extends Component<{
  match: match<OwnerQuery>;
}> {
  public render(): ReactElement {
    const { match: { params: { host, owner } } } = this.props;
    const hostArg = getHostTypeFromHost(host);

    return (
      <div>
        <div className="margin-sides-48">
          <h1 className="uk-heading-line uk-text-center pad-bot-48">
            <span>{`${host}/${owner}`}</span>
          </h1>

          <ReposQuery query={REPOS_QUERY} variables={{ owner, host: hostArg }}>
            {({ loading, error, data }): ReactElement => {
              if (loading) return <Loading />;

              // TODO: create prettier componenets
              if (error) {
                return (
                  <p>
                    Error :(
                    {error.message}
                  </p>
                );
              }

              if (!data || !data.repos) return <p>No data found</p>;
              const { repos } = data;

              return (
                <CardSet cards={repos.map((r): Card => ({
                  title: r.name,
                  body: 'Hello world',
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
