import React from 'react';
import {render} from '@testing-library/react';
import '../../node_modules/bootstrap/dist/css/bootstrap.css';
import '../index.css';
import { test, expect, vi } from 'vitest'
import { act } from 'react-dom/test-utils';
import NetworkCard from './NetworkCard';
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';

const sim:SimInfo = {"id":"27f06fe2-6e82-44b0-af4a-6975d169ff48","name":"test","description":"","steps":["/api/simulation/27f06fe2-6e82-44b0-af4a-6975d169ff48/step/72e1e5cb-3f31-4afd-818f-2293076547f7","/api/simulation/27f06fe2-6e82-44b0-af4a-6975d169ff48/step/f62bd8e5-2027-4fca-9e02-92c6ae6468ac"],"options":{"linkTeamPeers":true,"linkedTeamList":[],"evangelistList":[],"loneEvangelist":[],"initColors":[0],"maxColors":2,"agentsWithMemory":false}};
const steps:Array<Step> = [{"results":{"iterations":0,"colors":[[6,0]],"conversations":[0]},"id":"72e1e5cb-3f31-4afd-818f-2293076547f7","parent":"27f06fe2-6e82-44b0-af4a-6975d169ff48"},{"results":{"iterations":100,"colors":[[4,2],[6,0],[5,1],[4,2],[4,2],[4,2],[5,1],[5,1],[5,1],[5,1],[6,0],[4,2],[4,2],[4,2],[4,2],[5,1],[4,2],[4,2],[5,1],[5,1],[6,0],[5,1],[6,0],[6,0],[6,0],[6,0],[6,0],[5,1],[5,1],[5,1],[6,0],[4,2],[5,1],[4,2],[5,1],[5,1],[5,1],[6,0],[5,1],[4,2],[4,2],[6,0],[5,1],[5,1],[5,1],[6,0],[5,1],[5,1],[5,1],[5,1],[4,2],[5,1],[5,1],[5,1],[4,2],[5,1],[6,0],[6,0],[6,0],[6,0],[5,1],[5,1],[5,1],[5,1],[6,0],[4,2],[5,1],[6,0],[4,2],[5,1],[5,1],[5,1],[5,1],[4,2],[5,1],[4,2],[5,1],[4,2],[5,1],[6,0],[5,1],[5,1],[6,0],[6,0],[5,1],[5,1],[6,0],[4,2],[4,2],[5,1],[3,3],[3,3],[4,2],[3,3],[4,2],[4,2],[4,2],[6,0],[4,2],[4,2]],"conversations":[6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,6,6,6,6,6,5,6,6,6,6,6,6,6,5,6,6,6,6,6,5,6,5,6,6,6,6,5,6,6,6,6,6,6]},"id":"f62bd8e5-2027-4fca-9e02-92c6ae6468ac","parent":"27f06fe2-6e82-44b0-af4a-6975d169ff48"}]

vi.mock('../API/api');

test('renders without crashing', async () => {
    var result:any;
    await act(async () => {
        result = render(
            <NetworkCard sim={sim} steps={steps} readsim={function (id: string): void {
                throw new Error('Function not implemented.');
            } }/>
        );
    });
    expect(result).toBeDefined();
    expect(result.asFragment()).toMatchSnapshot();
});
