import React, { Component } from 'react';
import { match } from 'react-router-dom';
import { Location } from 'history';

import Nav from '../../components/Nav/Nav';

interface ProjectQuery {
  host: string;
  owner: string;
  name: string;
}

class Timeline extends Component<{
  match: match<ProjectQuery>;
  location: Location;
}> {
  render() {
    const { host, owner, name } = this.props.match.params;

    return (
      <div >
        <Nav location={location} />

        <div className="text-center title title-m pad-bot-m">
          {`${host}/${owner}/${name}`}
        </div>

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

export default Timeline;
