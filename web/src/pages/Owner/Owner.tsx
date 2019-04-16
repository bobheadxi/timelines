import React, { Component } from 'react';
import { match } from 'react-router-dom';
import { Location } from 'history';
import { Query } from "react-apollo";
import gql from "graphql-tag";

import Nav from '../../components/Nav/Nav';
import { getHostTypeFromHost } from '../../lib';

import { REPOS_QUERY } from '../../queries/repos';

interface OwnerQuery {
  host: string;
  owner: string;
}

class Owner extends Component<{
  match: match<OwnerQuery>;
  location: Location;
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

          <Query query={REPOS_QUERY} variables={{ owner, hostArg }}>
            {({ loading, error, data }) => {
              if (loading) return <p>Loading...</p>;
              if (error) {
                console.error(error);
                return <p>Error :( { error.message }</p>;
              }

              return <p>{JSON.stringify(data)}</p>
            }}
          </Query>
        </div>
      </div>
    );
  }
}

export default Owner;
