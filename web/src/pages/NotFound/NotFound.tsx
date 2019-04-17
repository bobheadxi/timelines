import React, { Component, ReactElement } from 'react';
import { RouteComponentProps } from 'react-router-dom';

class NotFound extends Component<RouteComponentProps> {
  public render(): ReactElement {
    return (
      <div>
        <div className="uk-position-center">
          <h1>Page Not Found</h1>
        </div>
      </div>
    );
  }
}

export default NotFound;
