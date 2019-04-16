import React, { Component } from 'react';
import { match } from 'react-router-dom';

import Nav from '../../components/Nav/Nav';
import { Burndown } from '../../components/vis';

interface RepoQuery {
  host: string;
  owner: string;
  name: string;
}

class Timeline extends Component<{
  match: match<RepoQuery>;
}> {
  render() {
    const { host, owner, name } = this.props.match.params;

    return (
      <div>
        <Nav location={location} />
        <div className="margin-sides-l">
          <h1 className="uk-heading-line uk-text-center pad-bot-l margin-sides-l">
            <span>{`${host}/${owner}/${name}`}</span>
          </h1>
        </div>
        <Burndown />
      </div>
    );
  }
}

export default Timeline;
