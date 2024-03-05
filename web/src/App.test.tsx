import React from 'react';
import {render, waitFor, screen} from '@testing-library/react';
import '../node_modules/bootstrap/dist/css/bootstrap.css';
import './index.css';
import {BrowserRouter} from 'react-router-dom';
import App from './App';
import { expect, test, vi } from 'vitest'
import { act } from 'react-dom/test-utils';

vi.mock('./API/api');

test('renders without crashing', async () => {
  var result:any;
  await act(async () => {
    result = render(
      <BrowserRouter>
          <App />
      </BrowserRouter>
    );
  });
  expect(result).toBeDefined();
  expect(result.asFragment()).toMatchSnapshot();
});

