import React from 'react';
import {Switch, Route} from 'react-router-dom';
import './App.css';
import Home from './Home/Home';
import Simulation from './Simulation/Simulation';

const Main = () => (
  <main>
    <Switch>
      <Route exact path='/' component={Home}/>
      <Route path='/simulation/:number' component={Simulation}/>
    </Switch>
  </main>
  )
  
  export default Main;