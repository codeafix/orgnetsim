import React from 'react';
import {render} from '@testing-library/react';
import '../../node_modules/bootstrap/dist/css/bootstrap.css';
import '../index.css';
import { test, vi, expect } from 'vitest'
import { act } from 'react-dom/test-utils';
import AgentColorChart from './AgentColorChart';
import { SimInfo } from '../API/SimInfo';

const sim:SimInfo = {"id":"27f06fe2-6e82-44b0-af4a-6975d169ff48","name":"test","description":"","steps":["/api/simulation/27f06fe2-6e82-44b0-af4a-6975d169ff48/step/72e1e5cb-3f31-4afd-818f-2293076547f7","/api/simulation/27f06fe2-6e82-44b0-af4a-6975d169ff48/step/f62bd8e5-2027-4fca-9e02-92c6ae6468ac"],"options":{"linkTeamPeers":true,"linkedTeamList":[],"evangelistList":[],"loneEvangelist":[],"initColors":[0],"maxColors":2,"agentsWithMemory":false}};

vi.mock('../API/api');

test('renders without crashing', async () => {
    var result:any;
    await act(async () => {
        result = render(
            <AgentColorChart sim={sim}/>
        );
    });
    expect(result).toBeDefined();
    expect(result.asFragment()).toMatchSnapshot();
});
