import React, { Component } from 'react';

class Loading extends Component {
  render() {
    return (
      <div className="fill-width fill-height">
        <div className="uk-position-center">
          <span uk-spinner="ratio: 4.5"></span>
        </div>
      </div>
    );
  }
}

export default Loading;
