import React, { Component, ReactElement } from 'react';
import { RouteComponentProps } from 'react-router-dom';

import surfer from 'assets/surfer@500.png';

class About extends Component<{
  api: string;
} & RouteComponentProps> {
  public render(): ReactElement {
    const { api } = this.props;

    // strip trailing path
    const apiRoot = api.substring(0, api.indexOf('/query'));

    // strip leading protocol for use with the shields.io badge
    const apiCheckTarget = apiRoot.replace(/(^\w+:|^)\/\//, '');

    return (
      <div className="margin-sides-48 uk-panel">
        <article className="uk-article">
          <img
            src={surfer}
            className="uk-align-right uk-margin-remove-adjacent uk-width-1-3"
            alt="timelines"
          />
          <p className="uk-text-lead">
            the stories of your projects and your communities.
          </p>
          <div>
            <a href={`${apiRoot}/playground`}>
              <img
                src={`https://img.shields.io/website/https/${apiCheckTarget}.svg?down_color=lightgrey&down_message=offline&label=api&up_message=online`}
                alt={`API Status (${api})`}
              />
            </a>
            &nbsp;
            <a href="https://github.com/bobheadxi/timelines">
              <img
                src="https://img.shields.io/github/last-commit/bobheadxi/timelines/master.svg?color=FC9514&label=last%20updated"
                alt="GitHub last commit (master)"
              />
            </a>
            &nbsp;
            <a href="https://join.slack.com/t/timelines-app/shared_invite/enQtNjEzMDE1NDk5NjAwLWZlN2ViZTE0NTNlNDZjZTNlOTNiNzZhZTZmNzgzZGVmNzcwZGE2NGJiN2QwNDQ0NzIyNmJlM2QzOTE4ZjQ3ZGE">
              <img
                src="https://img.shields.io/badge/slack-grey.svg?logo=slack"
                alt="Join the Slack community"
              />
            </a>
          </div>
          <p>Timelines is a web application! More details coming soon.</p>
        </article>

      </div>
    );
  }
}

export default About;
