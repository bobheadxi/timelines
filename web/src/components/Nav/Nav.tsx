import React, { Component } from 'react';

class Nav extends Component<{
  noTitles?: Boolean;
}> {
  render() {
    return (
      <nav
        className="uk-navbar-container uk-navbar-transparent margin-ends-l margin-sides-l"
        data-uk-navbar="mode: click">

        <div className="uk-navbar-left">
          {!this.props.noTitles 
            ? <div className="uk-navbar-item title title-m">
                <a href="/" className="uk-link-heading">Timelines</a>
              </div>
            : null}
        </div>

        <div className="uk-navbar-right">
          {!this.props.noTitles
            ? <div className="uk-navbar-item title title-m">
              <a href="/about" className="uk-link-heading">About</a>
            </div>
            : null}
        </div>
      </nav>
    );
  }
}

export default Nav;
