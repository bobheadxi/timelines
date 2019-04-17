import React, { Component, ReactElement } from 'react';
import { match } from 'react-router-dom';

interface HostQuery {
  host: string;
}

class Timeline extends Component<{
  match: match<HostQuery>;
}> {
  public render(): ReactElement {
    const { match: { params: { host } } } = this.props;

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
