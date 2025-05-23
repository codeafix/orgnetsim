import React, { useState, useEffect, useRef } from 'react';
import {Card, Table, Button, Modal, Form} from 'react-bootstrap';
import AgentColorChart from './AgentColorChart'
import API from '../API/api';
import Color from './Color';
import { SimInfo } from '../API/SimInfo';
import { Step } from '../API/Step';

type StepsCardProps = {
    sim:SimInfo;
    steps:Array<Step>;
    readsim(id:string): void;
}

const StepsCard = (props:StepsCardProps) => {
    const [showrunmodal, setShowrunmodal] = useState<boolean>(false);
    const [nostep, setnostep] = useState<boolean>(false);
    
    const [itcount, setitcount] = useState<number|undefined>(0);
    const [stepcount, setstepcount] = useState<number|undefined>(0);
    const [colors, setcolors] = useState<Array<number>>([]);
    const [filename, setfilename] = useState<string>("results.csv");

    const hlink = useRef<HTMLAnchorElement>(null);

    const setOptions = () => {
        setitcount(100);
        setstepcount(10);
    };

    useEffect(() => {
        if (!props.sim) {
            return
        }
        setOptions();
        setnostep((props.sim.steps || []).length === 0)
        const culrs = [];
        for (var i = 0; i < props.sim.options['maxColors']; i++) {
            culrs.push(i);
        }
        setcolors(culrs);
    },[props.sim, props.steps]);
 
    const handlerunshow = () => setShowrunmodal(true);

    const getresults = () => {
        API.getResultsCsv(props.sim).then(data => {
            setfilename(data.filename);
            const href = window.URL.createObjectURL(data.blob);
            const a = hlink.current;
            if(!a) return;
            a.href = href;
            a.click();
            a.href = '';
        }).catch(err => console.error(err));
    }

    const handlerunclose = () => {
        setShowrunmodal(false);
        setOptions();
    };

    const handlerun = () => {
        setShowrunmodal(false);
        const spec = {
            iterations: itcount||0,
            steps: stepcount||0
        };
        API.runsim(props.sim, spec).then(response => {
            props.readsim(response.parent);
        })
    }   

    return(
        <Card className="mb-2 mx-n2">
            <Card.Header>
                <Card.Title>Steps
                    <Button size="sm" className="btn btn-primary float-right" onClick={getresults} disabled={nostep}>Export Results</Button>
                    <a href="#/" ref={hlink} download={filename}>...</a>
                    <Button size="sm" className="btn btn-primary float-right mr-2" onClick={handlerunshow} disabled={nostep}>Run</Button>
                </Card.Title></Card.Header>
            <Card.Body className="small">
                <AgentColorChart sim={props.sim}/>
                <Table className="mb-n3" striped bordered size="sm">
                    <thead>
                        <tr>
                            <th key="h_iter">Iterations</th>
                            <th key="h_conv">Conversations</th>
                            {colors.map((color, i) => {
                                return <th key={"h_c_"+i}>{Color.colorFromVal(color)}</th>
                            })}
                        </tr>
                    </thead>
                    <tbody>
                        <StepsList steps={props.steps}/>
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
                    <Modal.Title>Run Simulation</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form.Group controlId="form-steps">
                        <Form.Label>Step count</Form.Label>
                        <Form.Control type="number" value={stepcount} onChange={(e:any) => setstepcount(e.target.valueAsNumber||"")}/>
                        <Form.Text className="text-muted">
                            The number of steps to run this simulation for
                        </Form.Text>
                    </Form.Group>
                    <Form.Group controlId="form-iterations">
                        <Form.Label>Iteration count</Form.Label>
                        <Form.Control type="number" value={itcount} onChange={(e:any) => setitcount(e.target.valueAsNumber||"")}/>
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

type StepsListProps = {
    steps:Array<Step>;
}


const StepsList = (props:StepsListProps) => {
    const [steps, setsteps] = useState<Array<Step>>([]);

    useEffect(() => {
        setsteps(props.steps || []);
    },[props.steps]);
    
    return steps.map(step => {        
        return(
            <StepItem key={step.id} step={step}/>
        );
    });
}

type StepItemProps = {
    step:Step;
}

const StepItem = (props:StepItemProps) => {
    const [iterations, setiterations] = useState<number>(0);
    const [conversations, setconversations] = useState<number>(0);
    const [colors, setcolors] = useState<Array<number>>([]);

    useEffect(() => {
        if(!props.step) return;
        const step = props.step;
        var itrs = step.results.iterations > 0 ? step.results.iterations - 1 : 0;
        setiterations(step.results.iterations);
        setconversations(step.results.conversations[itrs]);
        const culrs = [];
        for (var i = 0; i < step.results.colors[itrs].length; i++) {
            culrs.push(step.results.colors[itrs][i]);
        }
        setcolors(culrs);
    },[props.step, iterations]);
    
    return(
        <tr key={props.step.id}>
            <td key="iter">{iterations}</td>
            <td key="conv">{conversations}</td>
            {colors.map((colorcount, i) => {
                return <td key={"count_"+i}>{colorcount}</td>
            })}
        </tr>
    );
}

export default StepsCard