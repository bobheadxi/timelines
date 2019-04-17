import React, { Component } from 'react';

class Nav extends Component<{
  location: Location;
}> {
  render() {
    return (
      <nav
        className="uk-navbar-container uk-navbar-transparent margin-ends-l margin-sides-l"
        data-uk-navbar="mode: click">

        <div className="uk-navbar-left">
          {location.pathname == "/"
            ? null
            : <div className="uk-navbar-item title title-m">
                <a href="/" className="uk-link-heading">timelines</a>
              </div>}
        </div>

        <div className="uk-navbar-right">
          {location.pathname == "/about"
            ? null
            : <div className="uk-navbar-item title title-m">
                <a href="/about" className="uk-link-heading">about</a>
              </div>}
        </div>
      </nav>
    );
  }
}

export default Nav;
