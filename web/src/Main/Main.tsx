import React, { Component, Suspense, lazy } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from 'react-apollo';

import Loading from '../components/Loading/Loading';
import Nav from '../components/Nav/Nav';

const Home = lazy(() => import('../pages/Home/Home'));
const About = lazy(() => import('../pages/About/About'));
const Host = lazy(() => import('../pages/Host/Host'));
const Owner = lazy(() => import('../pages/Owner/Owner'));
const Timeline = lazy(() => import('../pages/Timeline/Timeline'));
const NotFound = lazy(() => import('../pages/NotFound/NotFound'));

const api = process.env.API_URL || 'https://timelines-api.herokuapp.com/query';
const client = new ApolloClient({
  uri: api,
})

class Main extends Component {
  render() {
    return (
      <ApolloProvider client={client}>
        <BrowserRouter>
          <div>
            <Nav location={location} />
            <Suspense fallback={<Loading />}>
              <Switch>
                <Route path="/" exact render={props => <Home {...props} />} />
                <Route path="/about" exact render={props => <About api={api} {...props} />} />
                <Route path="/:host" exact render={props => <Host {...props} />} />
                <Route path="/:host/:owner" exact render={props => <Owner {...props} />} />
                <Route path="/:host/:owner/:name" render={props => <Timeline {...props} />} />
                <Route path='*' exact render={props => <NotFound {...props} />} />
              </Switch>
            </Suspense>
          </div>
        </BrowserRouter>
      </ApolloProvider>
    );
  }
}

export default Main;
