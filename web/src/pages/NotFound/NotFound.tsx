import React, { Component } from 'react';
import { RouteComponentProps } from 'react-router-dom';

class NotFound extends Component<RouteComponentProps> {
  render() {
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
