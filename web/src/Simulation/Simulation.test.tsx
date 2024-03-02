import React from 'react';
import {render} from '@testing-library/react';
import '../../node_modules/bootstrap/dist/css/bootstrap.css';
import '../index.css';
import { test } from 'vitest'
import { vi } from 'vitest'
import { act } from 'react-dom/test-utils';
import Simulation from './Simulation';
import { BrowserRouter } from 'react-router-dom';

vi.mock('../API/api');

test('renders without crashing', async () => {
    await act(async () => {
        render(
            <BrowserRouter>
                <Simulation match={{
                        params: {
                            id: ''
                        }
                    }}/>
            </BrowserRouter>
        );
    });
});