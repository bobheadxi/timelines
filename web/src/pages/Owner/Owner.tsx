import React, { Component } from 'react';
import { match } from 'react-router-dom';
import { Location } from 'history';

import Nav from '../../components/Nav/Nav';
import Loading from '../../components/Loading/Loading';
import { getHostTypeFromHost } from '../../lib';

import { ReposQuery, REPOS_QUERY } from '../../lib/queries/repos';

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

          <ReposQuery query={REPOS_QUERY} variables={{ owner, hostArg }}>
            {({ loading, error, data }) => {
              if (loading) return <Loading />;

              // TODO: create prettier componenets
              if (error) {
                console.error(error);
                return <p>Error :( { error.message }</p>;
              }
              if (!data || !data.repos ) return <p>No data found</p>;

              // TODO: this style is somewhat often used, make a component to
              // do this
              return (
                <div className="margin-sides-xxl">
                  <div
                    className="uk-child-width-1-2@s uk-grid-match"
                    data-uk-scrollspy="target: > div; cls:uk-animation-fade; delay: 50"
                    data-uk-grid>
                    {data.repos.map(r => {
                      return (
                        <div>
                          <div className="uk-card uk-card-hover uk-card-default">
                            <div className="uk-card-body">
                              <h3 className="uk-card-title">{r.name}</h3>
                              <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
                              <a href={`${host}/${owner}/${name}`} className="uk-button uk-button-text">View project</a>
                            </div>
                          </div>
                        </div>
                      )
                    })}
                  </div>
                </div>
              )
            }}
          </ReposQuery>
        </div>
      </div>
    );
  }
}

export default Owner;
