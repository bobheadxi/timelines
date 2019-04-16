import React, { Component } from 'react';
import { match } from 'react-router-dom';

import Nav from '../../components/Nav/Nav';
import Loading from '../../components/Loading/Loading';
import CardSet from '../../components/CardSet/CardSet';
import { getHostTypeFromHost } from '../../lib';

import { ReposQuery, REPOS_QUERY } from '../../lib/queries/repos';

interface OwnerQuery {
  host: string;
  owner: string;
}

class Owner extends Component<{
  match: match<OwnerQuery>;
}> {
  render() {
    const { host, owner } = this.props.match.params;
    const hostArg = getHostTypeFromHost(host);

    return (
      <div>
        <Nav location={location} />
        <div className="margin-sides-l">
          <h1 className="uk-heading-line uk-text-center pad-bot-l margin-sides-l">
            <span>{`${host}/${owner}`}</span>
          </h1>

          <ReposQuery query={REPOS_QUERY} variables={{ owner, hostArg }}>
            {({ loading, error, data }) => {
              if (loading) return <Loading />;

              // TODO: create prettier componenets
              if (error) {
                console.error(error);
                return <p>Error :( { error.message }</p>;
              }

              if (!data || !data.repos) return <p>No data found</p>;
              const { repos } = data;

              return (
                <CardSet cards={repos.map(r => {
                  return {
                    title: r.name,
                    body: 'Hello world',
                    button: {
                      href: `${host}/${owner}/${name}`,
                      text: 'View Project',
                    },
                  }
                })} />
              )
            }}
          </ReposQuery>
        </div>
      </div>
    );
  }
}

export default Owner;
