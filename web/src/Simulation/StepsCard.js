import React, { useState, useEffect } from 'react';
import {Card, Table} from 'react-bootstrap';
import API from '../api';
import Color from './Color';

const StepsCard = (props) => {
    const [sim, setsim] = useState({steps: [], options:{}});
    const [colors, setcolors] = useState([]);

    useEffect(() => {
        setsim(props.sim);
        const culrs = [];
        for (var i = 0; i < props.sim.options['maxColors']; i++) {
            culrs.push(i);
        }
        setcolors(culrs);
    },[props.sim]);

    return(
        <Card className="mb-2 mx-n2">
            <Card.Header><Card.Title>Steps</Card.Title></Card.Header>
            <Card.Body className="small">
                <Table striped bordered hover>
                    <thead>
                        <tr>
                            <th>Iterations</th>
                            <th>Conversations</th>
                            {colors.map((color) => {
                                return <th>{Color.colorFromVal(color)}</th>
                            })}
                        </tr>
                    </thead>
                    <tbody>
                        <StepsList steps={sim.steps}/>
                    </tbody>
                </Table>
            </Card.Body>
        </Card>
    )
}

const StepsList = (props) => {
    const [steps, setsteps] = useState([]);

    useEffect(() => {
        setsteps(props.steps || []);
    },[props.steps]);
    
    return steps.map(step => {        
        return(
            <StepItem steppath={step}/>
        );
    });
}

const StepItem = (props) => {
    const [iterations, setiterations] = useState(0);
    const [conversations, setconversations] = useState(0);
    const [colors, setcolors] = useState([]);

    useEffect(() => {
        API.getStep(props.steppath).then(response => {
            const step = response;
            setiterations(step.results.iterations);
            setconversations(step.results.conversations[iterations]);
            const culrs = [];
            for (var i = 0; i < step.network['maxColors']; i++) {
                culrs.push(step.results.colors[iterations][i]);
            }
            setcolors(culrs);
        });
    },[props.steppath]);
    
    return(
        <tr>
            <td>{iterations}</td>
            <td>{conversations}</td>
            {colors.map((colorcount) => {
                return <td>{colorcount}</td>
            })}
        </tr>
    );
}

export default StepsCard