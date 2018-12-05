import './App.css';

import { Route, BrowserRouter as Router, Switch } from 'react-router-dom';
import React, {Component} from 'react';

import MainNavbar from './Components/Navbar'
import styled, { css } from 'react-emotion'
import Product from './Pages/Product'
import store from './store'
import { Provider } from 'react-redux'


class App extends Component {
  state = {
    isOpen: false
  }

  render() {
    return (
      <Provider store={store} >
        <Router>
          <div className="App">
            <MainNavbar />
            <Switch>
            <Route path="/products/:id" component={Product} />
            </Switch>
          </div>
        </Router>
      </Provider>
    );
  }
}

export default App;
