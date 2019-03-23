import React, { Component } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import Timeline from '../Timeline/Timeline';
import About from '../About/About';
import Home from '../Home/Home';

class Main extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <Switch>
            <Route path="/" exact component={Home} />
            <Route path="/about" component={About} />
            <Route path="/:project" component={Timeline} />
          </Switch>
        </div>
      </BrowserRouter>
    );
  }
}

export default Main;
