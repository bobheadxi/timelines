import React, { Component } from 'react';
import Nav from '../../components/Nav/Nav';


import banner from '../../assets/banner.png';

class Home extends Component {
  render() {
    return (
      <div >
        <Nav noTitles={true} />
        <div className="margin-sides-l">
          <header className="flex ai-center jc-center ">
            <img src={banner} alt="banner" height="100px" width="85%" />
          </header>

          <div
            className="uk-child-width-1-2@s uk-grid-match"
            data-uk-scrollspy="target: > div; cls:uk-animation-fade; delay: 50"
            data-uk-grid>
            <div>
              <div className="uk-card uk-card-hover uk-card-default">
                <div className="uk-card-body">
                  <h3 className="uk-card-title">About</h3>
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
                  <a href="/about" className="uk-button uk-button-text">Read more</a>
                </div>
              </div>
            </div>
            <div>
              <div className="uk-card uk-card-hover uk-card-default">
              <div className="uk-card-body">
                <h3 className="uk-card-title">Demo</h3>
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit.</p>
                  <a href="/github.com/bobheadxi/calories" className="uk-button uk-button-text">See the demo</a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default Home;
