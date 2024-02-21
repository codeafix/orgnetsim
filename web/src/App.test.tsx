import React from 'react';
import {render, screen} from '@testing-library/react';
import '../node_modules/bootstrap/dist/css/bootstrap.css';
import './index.css';
import {BrowserRouter} from 'react-router-dom';
import App from './App';
import { expect, test } from 'vitest'

test('renders without crashing', () => {
  render(
      <BrowserRouter>
          <App />
      </BrowserRouter>
  );
});
