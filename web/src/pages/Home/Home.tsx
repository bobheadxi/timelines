import React, { Component } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import Nav from '../../components/Nav/Nav';
import CardSet from '../../components/CardSet/CardSet';

import banner from '../../assets/banner.png';

class Home extends Component<{} & RouteComponentProps> {
  render() {
    return (
      <div >
        <div>
          <header className="flex ai-center jc-center ">
            <img src={banner} alt="banner" height="100px" width="85%" />
          </header>

          <div className="margin-sides-xxl">
            <CardSet cards={[
              {
                title: 'Demo User',
                body: 'Check out an example user or organization overview.',
                button: {
                  href: '/github.com/bobheadxi',
                  text: 'See the demo',
                },
              },
              {
                title: 'Demo Project',
                body: 'Check out an example project timeline.',
                button: {
                  href: '/github.com/bobheadxi/calories',
                  text: 'See the demo',
                },
              },
            ]}/>
          </div>
        </div>
      </div>
    );
  }
}

export default Home;
