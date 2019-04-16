import React, { Component } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Nav from '../../components/Nav/Nav';

class About extends Component<{} & RouteComponentProps> {
  render() {
    return (
      <div >
        <Nav location={location} />
        <div className="margin-sides-l">
          <div className="uk-child-width-1-2@s uk-grid-match " data-uk-grid>
            <div>
              <div className="uk-card uk-card-hover uk-card-default uk-card-body">
                  <h3 className="uk-card-title">Default</h3>
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
              </div>
            </div>
            <div>
              <div className="uk-card uk-card-hover uk-card-default uk-card-body">
                  <h3 className="uk-card-title">Primary</h3>
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default About;
