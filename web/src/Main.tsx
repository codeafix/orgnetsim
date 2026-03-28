import React from 'react';
import {Routes, Route} from 'react-router-dom';
import './App.css';
import Home from './Home/Home';
import Simulation from './Simulation/Simulation';

const Main = () => (
  <main>
    <Routes>
      <Route path='/' element={<Home/>}/>
      <Route path='/simulation/:id' element={<Simulation/>}/>
    </Routes>
  </main>
  )

  export default Main;