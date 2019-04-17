import React, { Component, ReactElement } from 'react';

class Loading extends Component {
  public render(): ReactElement {
    return (
      <div className="fill-width fill-height">
        <div className="uk-position-center">
          <span uk-spinner="ratio: 4.5" />
        </div>
      </div>
    );
  }
}

export default Loading;
