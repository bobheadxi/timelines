/* eslint-disable */

const UIkit: any = require('uikit');
const Icons: any = require('uikit/dist/js/uikit-icons');
UIkit.use(Icons);

import React from 'react';
import ReactDOM from 'react-dom';

import './styles/_all.scss';
import Main from './Main/Main';

import * as serviceWorker from './serviceWorker';

ReactDOM.render(<Main />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
