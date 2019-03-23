import React, { Component } from 'react';
import { match } from 'react-router-dom';
import { Location } from 'history';

import Nav from '../../components/Nav/Nav';

interface ProjectQuery {
  host: string;
  owner: string;
  name: string;
}

class Timeline extends Component<{
  match: match<ProjectQuery>;
  location: Location;
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
      </div>
    );
  }
}

export default Timeline;
