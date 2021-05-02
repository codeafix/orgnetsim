import React, { useState, useEffect } from 'react';
import {Card} from 'react-bootstrap';

const NetworkCard = (props) => {
    const [sim, setsim] = useState(props.sim);

    useEffect(() => {
        setsim(props.sim);
      },[props.sim]);

    return(
        <Card>
            <Card.Header><Card.Title>Network</Card.Title></Card.Header>
        </Card>
    )
}

export default NetworkCard