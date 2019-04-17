import React, { Component } from 'react';
import { match } from 'react-router-dom';

interface HostQuery {
  host: string;
}

class Timeline extends Component<{
  match: match<HostQuery>;
}> {
  render() {
    const { host } = this.props.match.params;

    return (
      <div>
        <div className="margin-sides-48">
          <h1 className="uk-heading-line uk-text-center pad-bot-48">
            <span>{`${host}`}</span>
          </h1>
        </div>
      </div>
    );
  }
}

export default Timeline;
