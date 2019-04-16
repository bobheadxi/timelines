import React, { Component } from 'react';
import { match } from 'react-router-dom';
import { Location } from 'history';

import Nav from '../../components/Nav/Nav';

interface HostQuery {
  host: string;
}

class Timeline extends Component<{
  match: match<HostQuery>;
  location: Location;
}> {
  render() {
    const { host } = this.props.match.params;

    return (
      <div>
        <Nav location={location} />
        <div className="margin-sides-l">
          <h1 className="uk-heading-line uk-text-center pad-bot-l margin-sides-l">
            <span>{`${host}`}</span>
          </h1>
        </div>
      </div>
    );
  }
}

export default Timeline;
