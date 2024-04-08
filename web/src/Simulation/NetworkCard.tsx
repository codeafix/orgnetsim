import React, { useState } from 'react';
import {Card, Button} from 'react-bootstrap'
import NetworkGraph from './NetworkGraph';
import ParseSettings from './ParseSettings';
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';

type NetworkCardProps = {
    sim:SimInfo;
    steps:Array<Step>;
    readsim(id:string): void;
}

const NetworkCard = (props:NetworkCardProps) => {
    const [showimpmodal, setshowimpmodal] = useState<boolean>(false);
    
    return(
        <Card className="mb-2 mx-n2">
            <Card.Header><Card.Title>Network<Button size="sm" className="btn btn-primary float-right" onClick={() => setshowimpmodal(true)}>Import</Button></Card.Title></Card.Header>
            <Card.Body>
                <NetworkGraph sim={props.sim} steps={props.steps}/>
            </Card.Body>
            <ParseSettings sim={props.sim} steps={props.steps} readsim={props.readsim} showimportmodal={showimpmodal} onclose={() => setshowimpmodal(false)}/>
        </Card>
    )
}

export default NetworkCard