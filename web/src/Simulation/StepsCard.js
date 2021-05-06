import React, { useState, useEffect } from 'react';
import {Card, Table, Button, Modal, Form} from 'react-bootstrap';
import API from '../api';
import Color from './Color';

const StepsCard = (props) => {
    const [sim, setsim] = useState({steps: [], options:{}});
    const [showrunmodal, setShowrunmodal] = useState(false);
    const [nostep, setnostep] = useState(false);
    
    const [itcount, setitcount] = useState(0);
    const [stepcount, setstepcount] = useState(0);
    const [colors, setcolors] = useState([]);

    const setOptions = () => {
        setitcount(100);
        setstepcount(10);
    };

    useEffect(() => {
        setsim(props.sim);
        setOptions();
        setnostep((props.sim.steps || []).length === 0)
        const culrs = [];
        for (var i = 0; i < props.sim.options['maxColors']; i++) {
            culrs.push(i);
        }
        setcolors(culrs);
    },[props.sim]);

 
    const handlerunshow = () => setShowrunmodal(true);

    const handlerunclose = () => {
        setShowrunmodal(false);
        setOptions();
    };

    const handlerun = () => {
        setShowrunmodal(false);
        const spec = {
            iterations: itcount,
            steps: stepcount
        };
        API.runsim(sim, spec).then(response => {
            setsim(response);
        }
        )
    }   

    return(
        <Card className="mb-2 mx-n2">
            <Card.Header>
                <Card.Title>Steps
                    <Button size="sm" className="btn btn-primary float-right">Export Results</Button>
                    <Button size="sm" className="btn btn-primary float-right mr-2" onClick={handlerunshow} disabled={nostep}>Run</Button>
                </Card.Title></Card.Header>
            <Card.Body className="small">
                <Table className="m-n3" striped bordered size="sm">
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
            <Modal
                show={showrunmodal}
                onHide={handlerunclose}
                backdrop="static"
                keyboard={false}
            >
                <Modal.Header closeButton>
                    <Modal.Title>Import Network</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form.Group controlId="form-steps">
                        <Form.Label>Step count</Form.Label>
                        <Form.Control type="number" value={stepcount} onChange={e => setstepcount(e.target.valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            The number of steps to run this simulation for
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-iterations">
                        <Form.Label>Iteration count</Form.Label>
                        <Form.Control type="number" value={itcount} onChange={e => setitcount(e.target.valueAsNumber)}/>
                        <Form.Text className="text-muted">
                            The number of iterations that will be computed in each step
                        </Form.Text>
                    </Form.Group>
                </Modal.Body>
                <Modal.Footer>
                    <Button variant="success" onClick={handlerun}>Run</Button>
                    <Button variant="secondary" onClick={handlerunclose}>Cancel</Button>
                </Modal.Footer>
            </Modal>

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