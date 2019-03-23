import React, { Component, Suspense, lazy } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import Loading from '../components/Loading/Loading';

const Home = lazy(() => import('../pages/Home/Home'));
const About = lazy(() => import('../pages/About/About'));
const Timeline = lazy(() => import('../pages/Timeline/Timeline'));
const NotFound = lazy(() => import('../pages/NotFound/NotFound'));

class Main extends Component {
  render() {
    return (
      <BrowserRouter>
        <div>
          <Suspense fallback={<Loading />}>
            <Switch>
              <Route path="/" exact render={props => <Home {...props} />} />
              <Route path="/about" render={props => <About {...props} />} />
              <Route path="/:host/:owner/:name" render={props => <Timeline {...props} />} />
              <Route path='*' exact render={props => <NotFound {...props} />} />
            </Switch>
          </Suspense>
        </div>
      </BrowserRouter>
    );
  }
}

export default Main;
