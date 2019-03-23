import React, { Component } from 'react';
import { Location } from 'history';

import Nav from '../../components/Nav/Nav';

class NotFound extends Component<{
  location: Location<any>;
}> {
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
