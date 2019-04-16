import React, { Component } from 'react';
import { Location } from 'history';

class Nav extends Component<{
  noTitles?: Boolean;
  location: any; // TODO: why can't this be of type Location?
}> {
  render() {
    const { noTitles } = this.props;
    return (
      <nav
        className="uk-navbar-container uk-navbar-transparent margin-ends-l margin-sides-l"
        data-uk-navbar="mode: click">

        <div className="uk-navbar-left">
          {noTitles
            ? null
            : <div className="uk-navbar-item title title-m">
                <a href="/" className="uk-link-heading">Timelines</a>
              </div>}
        </div>

        <div className="uk-navbar-right">
          {location.pathname == "/about"
            ? null
            : <div className="uk-navbar-item title title-m">
                <a href="/about" className="uk-link-heading">About</a>
              </div>}
        </div>
      </nav>
    );
  }
}

export default Nav;
