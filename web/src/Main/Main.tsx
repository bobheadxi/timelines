import React, {
  Component, Suspense, ReactElement, lazy,
} from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from 'react-apollo';

import Loading from '../components/Loading/Loading';
import Nav from '../components/Nav/Nav';

/* eslint-disable */
const Home = lazy(() => import('../pages/Home/Home'));
const About = lazy(() => import('../pages/About/About'));
const Host = lazy(() => import('../pages/Host/Host'));
const Owner = lazy(() => import('../pages/Owner/Owner'));
const Timeline = lazy(() => import('../pages/Timeline/Timeline'));
const NotFound = lazy(() => import('../pages/NotFound/NotFound'));
/* eslint-enable */

const api = process.env.API_URL || 'https://timelines-api.herokuapp.com/query';
const client = new ApolloClient({
  uri: api,
});

class Main extends Component {
  public render(): ReactElement {
    return (
      <ApolloProvider client={client}>
        <BrowserRouter>
          <div>
            <Nav location={window.location} />
            <Suspense fallback={<Loading />}>
              <Switch>
                <Route
                  path="/"
                  exact
                  render={(props): ReactElement => <Home {...props} />}
                />
                <Route
                  path="/about"
                  exact
                  render={(props): ReactElement => <About api={api} {...props} />}
                />
                <Route
                  path="/:host"
                  exact
                  render={(props): ReactElement => <Host {...props} />}
                />
                <Route
                  path="/:host/:owner"
                  exact
                  render={(props): ReactElement => <Owner {...props} />}
                />
                <Route
                  path="/:host/:owner/:name"
                  render={(props): ReactElement => <Timeline {...props} />}
                />
                <Route
                  path="*"
                  exact
                  render={(props): ReactElement => <NotFound {...props} />}
                />
              </Switch>
            </Suspense>
          </div>
        </BrowserRouter>
      </ApolloProvider>
    );
  }
}

export default Main;
