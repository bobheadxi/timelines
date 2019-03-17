import React, { Component } from 'react';
import banner from '../../assets/banner.png';

class Main extends Component {
  render() {
    // TODO: cards dont look right https://getuikit.com/docs/card
    return (
      <div>
        <header className="flex ai-center jc-center ">
          <img src={banner} alt="banner" height="100px" width="75%" />
        </header>

        <div className="uk-child-width-1-2@s uk-grid-match" uk-grid>
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
          <div>
            <div className="uk-card uk-card-hover uk-card-default uk-card-body">
                <h3 className="uk-card-title">Secondary</h3>
                <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Main;
