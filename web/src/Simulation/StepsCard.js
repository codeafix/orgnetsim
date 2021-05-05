import React, { useState, useEffect } from 'react';
import {Card} from 'react-bootstrap';

const StepsCard = (props) => {
    const [sim, setsim] = useState(props.sim);

    useEffect(() => {
        setsim(props.sim);
      },[props.sim]);

    return(
        <Card className="mb-2 mx-n2">
            <Card.Header><Card.Title>Steps</Card.Title></Card.Header>
        </Card>
    )
}

export default StepsCard