import React from 'react';
import ReactDOM from 'react-dom';
import './styles/_all.scss';
import Main from './components/Main/Main';
import * as serviceWorker from './serviceWorker';

const UIkit: any = require('uikit');
const Icons: any = require('uikit/dist/js/uikit-icons');

// load the icon plugin
UIkit.use(Icons);

ReactDOM.render(<Main />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
