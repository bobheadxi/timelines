import React, { Component } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import surfer from '../../assets/surfer@500.png';

class About extends Component<{
  api: string;
} & RouteComponentProps> {
  render() {
    const { api } = this.props;

    // strip leading protocol and trailing path for use to check API health
    const apiCheckTarget = api
      .substring(0, api.indexOf('/query'))
      .replace(/(^\w+:|^)\/\//, '');

    return (
      <div className="margin-sides-xxl uk-panel">
        <article className="uk-article">
          <img
            src={surfer}
            className="uk-align-right uk-margin-remove-adjacent uk-width-1-3" />
          <p className="uk-text-lead">
            the stories of your projects and your communities.
          </p>
          <div>
            <a href={api}>
              <img src={`https://img.shields.io/website/https/${apiCheckTarget}.svg?down_color=lightgrey&down_message=offline&label=api&up_message=online`}
                alt="API Status" />
            </a>
            &nbsp;
            <a href="https://github.com/bobheadxi">
              <img src="https://img.shields.io/github/last-commit/bobheadxi/timelines/master.svg?color=FC9514&label=last%20updated"
                alt="GitHub last commit (master)" />
            </a>
          </div>
          <p>Timelines is a web application.</p>
        </article>

      </div>
    );
  }
}

export default About;
