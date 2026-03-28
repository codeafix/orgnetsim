import React from 'react';
import {render} from '@testing-library/react';
import '../../node_modules/bootstrap/dist/css/bootstrap.css';
import '../index.css';
import { test, expect } from 'vitest'
import { vi } from 'vitest'
import { act } from 'react-dom/test-utils';
import Simulation from './Simulation';
import { MemoryRouter, Route, Routes } from 'react-router-dom';

vi.mock('../API/api');

test('renders without crashing', async () => {
    var result:any;
    await act(async () => {
        result = render(
            <MemoryRouter initialEntries={['/simulation/']}>
                <Routes>
                    <Route path='/simulation/:id' element={<Simulation/>}/>
                    <Route path='/simulation/' element={<Simulation/>}/>
                </Routes>
            </MemoryRouter>
        );
    });
    expect(result).toBeDefined();
    expect(result.asFragment()).toMatchSnapshot();
});