import React, { Component } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Nav from '../../components/Nav/Nav';

class NotFound extends Component<RouteComponentProps> {
  render() {
    return (
      <div>
        <Nav location={location} />
        <div className="uk-position-center">
          <h1>Page Not Found</h1>
        </div>
      </div>
    );
  }
}

export default NotFound;
