import React from 'react';
import {render} from '@testing-library/react';
import '../node_modules/bootstrap/dist/css/bootstrap.css';
import './index.css';
import {BrowserRouter} from 'react-router-dom';
import App from './App';
import { test, vi } from 'vitest'
import { act } from 'react-dom/test-utils';

vi.mock('./API/api');

test('renders without crashing', async () => {
  await act(async () => {
    render(
      <BrowserRouter>
          <App />
      </BrowserRouter>
    );
  });
});

