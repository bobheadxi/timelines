import React, { Component, ReactElement } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import CardSet from 'components/CardSet/CardSet';
import Contact from 'components/netlify-forms/Contact/Contact';
import { Error } from 'components/alerts';

import banner from 'assets/banner.png';

class Home extends Component<{} & RouteComponentProps> {
  public render(): ReactElement {
    return (
      <div>
        <div>
          <header className="flex ai-center jc-center ">
            <img src={banner} alt="banner" height="100px" width="85%" />
          </header>

          <div className="margin-sides-168">

            <Error message="WARNING: this project is a work in progress - nothing really works!" />

            <hr className="uk-divider-icon" />
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
            ]}
            />

            <hr className="uk-divider-icon" />

            {Contact}
          </div>
        </div>
      </div>
    );
  }
}

export default Home;
